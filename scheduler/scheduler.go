package scheduler

import (
	"fmt"
	"io/ioutil"
	"sync"
	"time"

	"github.com/blueskan/gopheart/log"
	"github.com/hako/durafmt"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/notifier"
	p "github.com/blueskan/gopheart/provider"
	"github.com/vmihailenco/msgpack"
)

type scheduler struct {
	providers       []p.Provider
	notifiers       map[string][]notifier.Notifier
	statistics      map[string]*p.Statistics
	persistence     bool
	persistenceLock sync.Mutex
}

func NewScheduler(
	providers []p.Provider,
	notifiers map[string][]notifier.Notifier,
	statistics map[string]*p.Statistics,
	persistOnDisk bool,
) p.Scheduler {
	scheduler := &scheduler{
		providers:   providers,
		notifiers:   notifiers,
		persistence: persistOnDisk,
	}

	if statistics != nil {
		scheduler.statistics = statistics
	} else {
		scheduler.statistics = make(map[string]*p.Statistics)
	}

	for _, provider := range providers {
		name := provider.GetName()

		if _, exists := scheduler.statistics[name]; !exists {
			initialAuditLogs := make([]p.AuditLog, 0)

			initialAuditLogs = append(initialAuditLogs, p.AuditLog{
				Timestamp:     time.Now(),
				PreviousState: p.Unknown,
				NewState:      p.Unknown,
			})

			scheduler.statistics[name] = &p.Statistics{
				ServiceName:           name,
				CurrentFailureCount:   0,
				CurrentSuccessCount:   0,
				CurrentErrorThreshold: 0,
				RunningInterval:       provider.GetInterval(),
				NextRunAt:             time.Now().Add(provider.GetInterval()),
				State:                 p.Unknown,
				AuditLogs:             initialAuditLogs,
			}
		}
	}

	return scheduler
}

func (s *scheduler) GetStatistics() map[string]*p.Statistics {
	return s.statistics
}

func (s *scheduler) Schedule() {
	go func() {
		log.Info(fmt.Sprintf("Starting scheduling healthchecks with %d providers", len(s.providers)))

		for _, provider := range s.providers {
			go func(provider p.Provider) {
				for {
					name := provider.GetName()
					log.Info(fmt.Sprintf("Healthcheck scheduling starting for `%s`.", name))

					statistics := s.statistics[name]

					s.canBeRunning(statistics)

					log.Info(fmt.Sprintf("Doing healthcheck for service %s.", name))

					resStartTime := time.Now()
					isHealthy := provider.Heartbeat()
					elapsedTime := time.Now().Sub(resStartTime)

					s.evaluateHealthCheckResponse(isHealthy, statistics, resStartTime, elapsedTime, name)

					previousState := statistics.State

					s.changeServiceState(isHealthy, statistics, provider)
					isStateChanged := previousState != statistics.State

					if isStateChanged {
						statistics.AuditLogs = append([]p.AuditLog{
							{
								Timestamp:     time.Now(),
								PreviousState: previousState,
								NewState:      statistics.State,
							},
						}, statistics.AuditLogs...)
					}

					s.notify(name, statistics, isStateChanged)

					if s.persistence {
						s.persistStats()
					}
				}
			}(provider)
		}
	}()
}

func (s *scheduler) changeServiceState(isHealthy bool, statistics *p.Statistics, provider p.Provider) {
	statistics.LastRunAt = time.Now()
	statistics.NextRunAt = time.Now().Add(statistics.RunningInterval)

	if isHealthy {
		statistics.CurrentSuccessCount++

		if statistics.CurrentSuccessCount >= provider.GetUpThreshold() {
			statistics.CurrentFailureCount = 0
			statistics.State = p.Healthy
		} else if statistics.State != p.Unknown {
			statistics.State = p.Sick
		}
	} else if statistics.CurrentFailureCount <= provider.GetDownThreshold() {
		statistics.CurrentFailureCount++

		statistics.State = p.Sick
	} else {
		statistics.CurrentFailureCount++
		statistics.CurrentSuccessCount = 0

		statistics.State = p.UnHealthy
	}
}

func (s *scheduler) evaluateHealthCheckResponse(isHealthy bool, statistics *p.Statistics, resStartTime time.Time, elapsedTime time.Duration, name string) {
	if isHealthy {
		statistics.ResponseLogs = append([]p.ResponseLog{
			{
				Timestamp:   resStartTime,
				ElapsedTime: elapsedTime,
			},
		}, statistics.ResponseLogs...)

		log.Success(fmt.Sprintf(
			"Healthcheck successful for `%s`, elapsed time: %s",
			name,
			durafmt.Parse(elapsedTime).String(),
		))
	} else {
		log.Error(fmt.Sprintf(
			"Healthcheck failed for `%s`, elapsed time: %s",
			name,
			durafmt.Parse(elapsedTime).String(),
		))
	}
}

func (s *scheduler) persistStats() {
	s.persistenceLock.Lock()

	log.Info(fmt.Sprintf("Starting persist healthcheck statistics data to disk"))

	statistics := s.GetStatistics()
	data, _ := msgpack.Marshal(&statistics)

	err := ioutil.WriteFile(config.DbPath, data, 0644)

	if err != nil {
		log.Error(fmt.Sprintf("Persistence operating failed"))
		panic(err)
	}

	log.Success(fmt.Sprintf("Persistence operating completed successfully"))

	s.persistenceLock.Unlock()
}

func (s *scheduler) canBeRunning(statistics *p.Statistics) {
	for {
		canBeRunning := time.Now().After(statistics.NextRunAt) || statistics.LastRunAt.Unix() <= 0

		if canBeRunning {
			break
		}

		duration := statistics.NextRunAt.Sub(time.Now())

		time.Sleep(duration)
	}
}

func (s *scheduler) notify(name string, statistics *p.Statistics, isStateChanged bool) {
	go func() {
		// TODO fix this wrong usage, we dont need to access first element though
		threshold := s.notifiers[name][0].GetThreshold()

		latestAuditLog := statistics.AuditLogs[0]

		// Notifying rules
		if latestAuditLog.PreviousState == p.Unknown && latestAuditLog.NewState == p.Healthy {
			return
		}

		if !isStateChanged && statistics.State == p.Healthy {
			return
		}

		if !isStateChanged && (statistics.State != p.Healthy) {
			if statistics.State != p.Unknown {
				statistics.CurrentErrorThreshold++
			}

			if statistics.CurrentErrorThreshold <= threshold {
				return
			}
		}

		statistics.CurrentErrorThreshold = 0

		for _, nfier := range s.notifiers[name] {
			log.Info(fmt.Sprintf("Notify service `%s` status change to `%s`.", name, nfier.GetName()))

			err := nfier.Notify(*statistics)

			if err != nil {
				log.Error(fmt.Sprintf("Cannot notify service `%s` status change to `%s`.", name, nfier.GetName()))
			} else {
				log.Success(fmt.Sprintf("Successfully notify service `%s` status change to `%s`.", name, nfier.GetName()))
			}
		}
	}()
}

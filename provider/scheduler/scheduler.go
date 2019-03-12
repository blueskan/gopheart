package scheduler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/notifier"
	p "github.com/blueskan/gopheart/provider"
	"github.com/vmihailenco/msgpack"
)

type scheduler struct {
	providers   []p.Provider
	notifiers   map[string][]notifier.Notifier
	statistics  map[string]*p.Statistics
	persistence bool
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

	// if config change we need doing some jobs here
	if statistics != nil {
		scheduler.statistics = statistics

		return scheduler
	}

	scheduler.statistics = make(map[string]*p.Statistics)

	for _, provider := range providers {
		name := provider.GetName()

		if _, exists := scheduler.statistics[name]; !exists {
			scheduler.statistics[name] = &p.Statistics{
				ServiceName:         name,
				CurrentFailureCount: 0,
				CurrentSuccessCount: provider.GetUpThreshold(),
				RunningInterval:     provider.GetInterval(),
				NextRunAt:           time.Now().Add(provider.GetInterval()),
				State:               p.Healthy,
				AuditLogs:           make([]p.AuditLog, 0),
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
		log.Printf("[INFO] Starting scheduling with %d providers.", len(s.providers))

		for _, provider := range s.providers {
			go func(provider p.Provider) {
				for {
					name := provider.GetName()
					log.Printf("[INFO] Healthcheck scheduling starting for `%s`.", name)

					statistics := s.statistics[name]

					for {
						canBeRunning := time.Now().After(statistics.NextRunAt) || statistics.LastRunAt.Unix() <= 0

						if canBeRunning {
							break
						}

						duration := statistics.NextRunAt.Sub(time.Now())

						log.Printf("[INFO] Healthcheck scheduling waiting %f seconds for `%s`.", duration.Seconds(), name)
						time.Sleep(duration)
					}

					log.Printf("[INFO] Doing healthcheck for service %s.", name)
					resStartTime := time.Now()
					isHealthy := provider.Heartbeat()
					elapsedTime := time.Now().Sub(resStartTime).Nanoseconds()

					status := "failed"
					if isHealthy {
						statistics.ResponseLogs = append([]p.ResponseLog{
							{
								Timestamp:   resStartTime,
								ElapsedTime: elapsedTime,
							},
						}, statistics.ResponseLogs...)

						status = "success"
					}

					log.Printf("[INFO] Completed healthcheck for service %s, Status: %s.", name, status)

					statistics.LastRunAt = time.Now()
					statistics.NextRunAt = time.Now().Add(statistics.RunningInterval)

					previousState := statistics.State

					// Please refactor this area: extract to function etc..
					if isHealthy {
						statistics.CurrentSuccessCount++

						if statistics.CurrentSuccessCount >= provider.GetUpThreshold() {
							statistics.CurrentFailureCount = 0
							statistics.State = p.Healthy
						} else {
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

					if previousState != statistics.State {
						statistics.AuditLogs = append([]p.AuditLog{
							{
								Timestamp:     time.Now(),
								PreviousState: previousState,
								NewState:      statistics.State,
							},
						}, statistics.AuditLogs...)

						// Notify with notifiers
						for _, notifier := range s.notifiers[name] {
							notifier.Notify(*statistics)
						}
					}

					json, _ := json.Marshal(statistics)
					log.Printf("[INFO] Latest status for service `%s`:\n\n%s\n\n", name, string(json))

					if s.persistence {
						log.Printf("[INFO] Starting persist on a disk")

						statistics := s.GetStatistics()
						data, _ := msgpack.Marshal(&statistics)

						ioutil.WriteFile(config.DbPath, data, 0644)

						log.Printf("[INFO] Finish persist on a disk")
					}
				}
			}(provider)
		}
	}()
}

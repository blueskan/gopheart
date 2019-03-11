package provider

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type scheduler struct {
	providers  []Provider
	statistics map[string]*Statistics
}

func NewScheduler(providers []Provider) Scheduler {
	scheduler := &scheduler{
		providers: providers,
	}
	scheduler.statistics = make(map[string]*Statistics)

	for _, provider := range providers {
		name := provider.GetName()

		// TODO: Support persistence for statistics

		if _, exists := scheduler.statistics[name]; !exists {
			scheduler.statistics[name] = &Statistics{
				CurrentFailureCount: 0,
				CurrentSuccessCount: provider.GetUpThreshold(),
				RunningInterval:     provider.GetInterval(),
				NextRunAt:           time.Now().Add(provider.GetInterval()),
				State:               Healthy,
			}
		}
	}

	return scheduler
}

func (s *scheduler) GetStatistics() map[string]*Statistics {
	return s.statistics
}

func (s *scheduler) Schedule() {
	var waitGroup sync.WaitGroup

	for {
		waitGroup.Add(len(s.providers))

		log.Printf("[INFO] Starting scheduling with %d providers.", len(s.providers))

		for _, provider := range s.providers {
			go func(provider Provider) {
				name := provider.GetName()
				log.Printf("[INFO] Healthcheck scheduling starting for `%s`.", name)

				statistics := s.statistics[name]

				waitCh := make(chan bool)

				go func() {
					for {
						canBeRunning := time.Now().After(statistics.NextRunAt) || statistics.LastRunAt.Unix() <= 0

						if canBeRunning {
							break
						}

						duration := statistics.NextRunAt.Sub(time.Now())

						log.Printf("[INFO] Healthcheck scheduling waiting %f seconds for `%s`.", duration.Seconds(), name)
						time.Sleep(duration)
					}

					waitCh <- true
				}()

				<-waitCh

				log.Printf("[INFO] Doing healthcheck for service %s.", name)
				isHealthy := provider.Heartbeat()

				status := "failed"
				if isHealthy {
					status = "success"
				}

				log.Printf("[INFO] Completed healthcheck for service %s, Status: %s.", name, status)

				statistics.LastRunAt = time.Now()
				statistics.NextRunAt = time.Now().Add(statistics.RunningInterval)

				if isHealthy {
					statistics.CurrentSuccessCount++

					if statistics.CurrentSuccessCount > 3 {
						statistics.CurrentFailureCount = 0
						statistics.State = Healthy
					} else {
						statistics.State = Sick
					}
				} else if statistics.CurrentFailureCount <= 3 {
					statistics.CurrentFailureCount++

					statistics.State = Sick
				} else {
					statistics.CurrentFailureCount++

					statistics.State = UnHealthy
				}

				json, _ := json.Marshal(statistics)
				log.Printf("[INFO] Latest status for service `%s`:\n\n%s\n\n", name, string(json))

				waitGroup.Done()
			}(provider)
		}

		waitGroup.Wait()
	}
}

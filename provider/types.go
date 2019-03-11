package provider

import (
	"time"
)

type Provider interface {
	GetName() string
	GetInterval() time.Duration
	GetDownThreshold() int
	GetUpThreshold() int
	Heartbeat() bool
}

type Scheduler interface {
	GetStatistics() map[string]*Statistics
	Schedule()
}

type StateCode int8

const (
	UnHealthy StateCode = 0
	Sick      StateCode = 1
	Healthy   StateCode = 2
)

type Statistics struct {
	RunningInterval     time.Duration
	LastRunAt           time.Time
	NextRunAt           time.Time
	CurrentFailureCount int
	CurrentSuccessCount int
	State               StateCode
}

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

type StateCode string

const (
	UnHealthy StateCode = "Unhealthy"
	Sick      StateCode = "Sick"
	Healthy   StateCode = "Healthy"
)

type ResponseLog struct {
	Timestamp   time.Time `json:"timestamp"`
	ElapsedTime int64     `json:"elapsed_time"`
}

type AuditLog struct {
	Timestamp     time.Time `json:"timestamp"`
	PreviousState StateCode `json:"previous_state"`
	NewState      StateCode `json:"new_state"`
}

type Statistics struct {
	ServiceName         string        `json:"service_name"`
	RunningInterval     time.Duration `json:"running_interval"`
	LastRunAt           time.Time     `json:"last_run_at"`
	NextRunAt           time.Time     `json:"next_run_at"`
	CurrentFailureCount int           `json:"current_failure_count"`
	CurrentSuccessCount int           `json:"current_success_count"`
	State               StateCode     `json:"state"`
	AuditLogs           []AuditLog    `json:"audit_logs"`
	ResponseLogs        []ResponseLog `json:"response_logs"`
}

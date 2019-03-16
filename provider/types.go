package provider

import (
	"github.com/hako/durafmt"
	"time"
)

type Provider interface {
	GetName() string
	GetInterval() time.Duration
	GetDownThreshold() int64
	GetUpThreshold() int64
	Heartbeat() bool
}

type Scheduler interface {
	GetStatistics() map[string]*Statistics
	Schedule()
}

type StateCode string

const (
	Unknown   StateCode = "Unknown"
	UnHealthy StateCode = "Unhealthy"
	Sick      StateCode = "Sick"
	Healthy   StateCode = "Healthy"
)

type ResponseLog struct {
	Timestamp   time.Time     `json:"timestamp"`
	ElapsedTime time.Duration `json:"elapsed_time"`
}

type ResponseLogRepresentation struct {
	Timestamp   time.Time     `json:"timestamp"`
	ElapsedTime string `json:"elapsed_time"`
}

func (rl ResponseLog) Representation() ResponseLogRepresentation {
	return ResponseLogRepresentation{
		Timestamp: rl.Timestamp,
		ElapsedTime: durafmt.Parse(rl.ElapsedTime).String(),
	}
}

type AuditLog struct {
	Timestamp     time.Time `json:"timestamp"`
	PreviousState StateCode `json:"previous_state"`
	NewState      StateCode `json:"new_state"`
}

type Statistics struct {
	ServiceName           string        `json:"service_name"`
	RunningInterval       time.Duration `json:"running_interval"`
	LastRunAt             time.Time     `json:"last_run_at"`
	NextRunAt             time.Time     `json:"next_run_at"`
	CurrentFailureCount   int64         `json:"current_failure_count"`
	CurrentSuccessCount   int64         `json:"current_success_count"`
	CurrentErrorThreshold int           `json:"current_error_threshold"`
	State                 StateCode     `json:"state"`
	AuditLogs             []AuditLog    `json:"audit_logs"`
	ResponseLogs          []ResponseLog `json:"response_logs"`
}

type StatisticsRepresentation struct {
	State StateCode `json:"state"`
	RunningInterval string `json:"running_interval"`
	LastRunAt             time.Time     `json:"last_run_at"`
	NextRunAt             time.Time     `json:"next_run_at"`
	AuditLogs []AuditLog `json:"audit_logs"`
	ResponseLogs          []ResponseLogRepresentation `json:"response_logs"`
}

func (s Statistics) Representation(auditLogLimit, responseLogLimit int) StatisticsRepresentation {
	rlRepresentation := make([]ResponseLogRepresentation, 0)

	for _, val := range s.ResponseLogs {
		rlRepresentation = append(rlRepresentation, val.Representation())
	}

	aLogLimit := auditLogLimit

	if len(s.AuditLogs) < auditLogLimit {
		aLogLimit = len(s.AuditLogs)
	}

	rLogLimit := responseLogLimit

	if len(s.ResponseLogs) < responseLogLimit {
		rLogLimit = len(s.ResponseLogs)
	}

	return StatisticsRepresentation{
		State: s.State,
		RunningInterval: durafmt.Parse(s.RunningInterval).String(),
		LastRunAt: s.LastRunAt,
		NextRunAt: s.NextRunAt,
		AuditLogs: s.AuditLogs[:aLogLimit],
		ResponseLogs: rlRepresentation[:rLogLimit],
	}
}

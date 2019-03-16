package notifier

import (
	"github.com/hako/durafmt"
	"strings"
	"time"

	"github.com/blueskan/gopheart/provider"
)

func ComposeMessage(template string, statistics provider.Statistics) string {
	latestAuditLog := statistics.AuditLogs[1]

	output := template

	duration := time.Now().Sub(latestAuditLog.Timestamp)

	output = strings.Replace(output, "${service_name}", statistics.ServiceName, -1)
	output = strings.Replace(output, "${duration}", durafmt.Parse(duration).String(), -1)
	output = strings.Replace(output, "${previous_state}", string(latestAuditLog.PreviousState), -1)
	output = strings.Replace(output, "${new_state}", string(latestAuditLog.NewState), -1)

	return output
}
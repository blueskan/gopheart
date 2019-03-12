package notifier

import (
	"fmt"
	"strings"
	"time"

	"github.com/blueskan/gopheart/provider"
)

func ComposeMessage(template string, statistics provider.Statistics) string {
	latestAuditLog := statistics.AuditLogs[0]

	output := template

	duration := time.Now().Sub(latestAuditLog.Timestamp)

	output = strings.Replace(output, "${service_name}", statistics.ServiceName, -1)
	output = strings.Replace(output, "${duration}", fmtDuration(duration), -1)
	output = strings.Replace(output, "${previous_state}", string(latestAuditLog.PreviousState), -1)
	output = strings.Replace(output, "${new_state}", string(latestAuditLog.NewState), -1)

	return output
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02d:%02d", h, m)
}

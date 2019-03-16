package custom

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/blueskan/gopheart/provider"
)

type customProvider struct {
	name          string
	command       string
	timeout       time.Duration
	interval      time.Duration
	downThreshold int64
	upThreshold   int64
}

func NewCustomProvider(
	name, command string,
	timeout, interval time.Duration,
	downThreshold, upThreshold int64,
) provider.Provider {
	return &customProvider{
		name:          name,
		command:       command,
		timeout:       timeout,
		interval:      interval,
		downThreshold: downThreshold,
		upThreshold:   upThreshold,
	}
}

func (cp customProvider) GetName() string {
	return cp.name
}

func (cp customProvider) GetInterval() time.Duration {
	return cp.interval
}

func (cp customProvider) GetDownThreshold() int64 {
	return cp.downThreshold
}

func (cp customProvider) GetUpThreshold() int64 {
	return cp.upThreshold
}

func (cp customProvider) Heartbeat() bool {
	command := strings.Split(cp.command, " ")

	cmd := exec.Command(command[0], command[1:]...)
	env := os.Environ()
	env = append(env, fmt.Sprintf("NAME=%s", cp.name))
	env = append(env, fmt.Sprintf("TIMEOUT=%s", cp.timeout))
	cmd.Env = env

	startErr := cmd.Start()

	if startErr != nil {
		return false
	}

	err := cmd.Wait()

	if err != nil {
		return false
	}

	return true
}

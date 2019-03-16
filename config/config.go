package config

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var ConfigPath = "./config.yaml"
var DbPath = "./gopheart.db"

type NotifierService struct {
	Url     string `yaml:"url" json:"url"`
	Message string `yaml:"message" json:"message"`
	Extra   string `yaml:"extra" json:"extra"`
}

type Notifiers struct {
	Threshold int                        `yaml:"threshold" json:"threshold"`
	Services  map[string]NotifierService `yaml:"services" json:"services"`
}

type RetryPolicy struct {
	Timeout       string `yaml:"timeout" json:"timeout"`
	DownThreshold int64  `yaml:"down_threshold" json:"down_threshold"`
	UpThreshold   int64  `yaml:"up_threshold" json:"up_threshold"`
}

type HealthCheck struct {
	Type          string      `yaml:"type" json:"type"`
	Source        string      `yaml:"source" json:"source"`
	CheckInterval string      `yaml:"check_interval" json:"check_interval"`
	RetryPolicy   RetryPolicy `yaml:"retry_policy" json:"retry_policy"`
	Notifiers     Notifiers   `yaml:"notifiers" json:"notifiers"`
}

type WebUI struct {
	Port              string `yaml:"port" json:"port"`
	FailureStatusCode string `yaml:"failure_status_code" json:"failure_status_code"`
	AuditLogLimit     int    `yaml:"audit_log_limit" json:"audit_log_limit"`
	ResponseLogLimit  int    `yaml:"response_log_limit" json:"response_log_limit"`
}

type GlobalConfiguration struct {
	WebUI                      WebUI  `yaml:"web_ui" json:"web_ui"`
	CollectStats               bool   `yaml:"collect_stats" json:"collect_stats"`
	AuditLogRotation           int    `yaml:"audit_log_rotation" json:"audit_log_rotation"`
	AuditLogRotationEnabled    bool   `yaml:"audit_log_rotation_enabled" json:"audit_log_rotation_enabled"`
	ResponseLogRotation        int    `yaml:"response_log_rotation" json:"response_log_rotation"`
	ResponseLogRotationEnabled bool   `yaml:"response_log_rotation_enabled" json:"response_log_rotation_enabled"`
	CheckInterval              string `yaml:"check_interval" json:"check_interval"`
	Notifiers                  Notifiers
	RetryPolicy                RetryPolicy
}

type Config struct {
	Global       GlobalConfiguration    `yaml:"global" json:"global"`
	HealthChecks map[string]HealthCheck `yaml:"health_checks" json:"health_checks"`
}

func (c *Config) FromJson(path string) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	json.Unmarshal(configFile, c)
	c.overrideOmitValuesWithDefaults()
}

func (c *Config) FromYaml(path string) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	yaml.Unmarshal(configFile, c)
	c.overrideOmitValuesWithDefaults()
}

func (c *Config) overrideOmitValuesWithDefaults() {
	for key, value := range c.HealthChecks {
		for serviceName, service := range value.Notifiers.Services {
			if len(service.Url) <= 0 {
				service.Url = c.Global.Notifiers.Services[serviceName].Url
			}

			if len(service.Message) <= 0 {
				service.Message = c.Global.Notifiers.Services[serviceName].Message
			}

			if len(service.Extra) <= 0 {
				service.Extra = c.Global.Notifiers.Services[serviceName].Extra
			}

			c.HealthChecks[key].Notifiers.Services[serviceName] = service
		}

		if value.Notifiers.Threshold <= 0 {
			value.Notifiers.Threshold = c.Global.Notifiers.Threshold
		}

		if len(value.CheckInterval) <= 0 {
			value.CheckInterval = c.Global.CheckInterval
		}

		if len(value.RetryPolicy.Timeout) <= 0 {
			value.RetryPolicy.Timeout = c.Global.RetryPolicy.Timeout
		}

		if value.RetryPolicy.DownThreshold <= 0 {
			value.RetryPolicy.DownThreshold = c.Global.RetryPolicy.DownThreshold
		}

		if value.RetryPolicy.UpThreshold <= 0 {
			value.RetryPolicy.UpThreshold = c.Global.RetryPolicy.UpThreshold
		}

		c.HealthChecks[key] = value
	}
}

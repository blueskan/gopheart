package main

import (
	"io/ioutil"
	"strconv"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/http"
	"github.com/blueskan/gopheart/provider"
	"github.com/blueskan/gopheart/provider/factory"
	"github.com/vmihailenco/msgpack"
)

func main() {
	configuration := config.Config{}
	configuration.FromYaml(config.ConfigPath)

	providerFactory := factory.NewProviderFactory()
	providers := make([]provider.Provider, 0)

	for key, value := range configuration.HealthChecks {
		providers = append(providers, providerFactory.CreateProvider(key, value))
	}

	var statistics map[string]*provider.Statistics

	if configuration.Global.CollectStats {
		database, err := ioutil.ReadFile(config.DbPath)
		if err == nil {
			msgpack.Unmarshal(database, &statistics)
		}
	}

	scheduler := provider.NewScheduler(providers, statistics, configuration.Global.CollectStats)
	scheduler.Schedule()

	failureStatusCode, _ := strconv.Atoi(configuration.Global.WebUI.FailureStatusCode)
	httpServer := http.NewHttpServer(scheduler, failureStatusCode)
	httpServer.Listen(configuration.Global.WebUI.Port)
}

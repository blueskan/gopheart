package main

import (
	"io/ioutil"
	"strconv"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/http"
	"github.com/blueskan/gopheart/notifier"
	nFactory "github.com/blueskan/gopheart/notifier/factory"
	"github.com/blueskan/gopheart/provider"
	pFactory "github.com/blueskan/gopheart/provider/factory"
	s "github.com/blueskan/gopheart/provider/scheduler"
	"github.com/vmihailenco/msgpack"
)

func main() {
	configuration := config.Config{}
	configuration.FromYaml(config.ConfigPath)

	providerFactory := pFactory.NewProviderFactory()
	notifierFactory := nFactory.NewProviderFactory()

	providers := make([]provider.Provider, 0)
	notifiers := make(map[string][]notifier.Notifier)

	for key, value := range configuration.HealthChecks {
		providers = append(providers, providerFactory.CreateProvider(key, value))

		notifiers[key] = make([]notifier.Notifier, 0)

		for service, notifier := range value.Notifiers.Services {
			notifiers[key] = append(notifiers[key], notifierFactory.CreateNotifier(key, service, notifier))
		}
	}

	var statistics map[string]*provider.Statistics

	if configuration.Global.CollectStats {
		database, err := ioutil.ReadFile(config.DbPath)
		if err == nil {
			msgpack.Unmarshal(database, &statistics)
		}
	}

	scheduler := s.NewScheduler(
		providers,
		notifiers,
		statistics,
		configuration.Global.CollectStats,
	)
	scheduler.Schedule()

	failureStatusCode, _ := strconv.Atoi(configuration.Global.WebUI.FailureStatusCode)
	httpServer := http.NewHttpServer(scheduler, failureStatusCode)
	httpServer.Listen(configuration.Global.WebUI.Port)
}

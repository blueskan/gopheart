package main

import (
	"flag"
	"fmt"
	"github.com/blueskan/gopheart/log"
	"io/ioutil"
	"path"
	"strconv"

	"github.com/blueskan/gopheart/config"
	"github.com/blueskan/gopheart/http"
	"github.com/blueskan/gopheart/notifier"
	nFactory "github.com/blueskan/gopheart/notifier/factory"
	"github.com/blueskan/gopheart/provider"
	pFactory "github.com/blueskan/gopheart/provider/factory"
	s "github.com/blueskan/gopheart/scheduler"
	_ "github.com/dimiro1/banner/autoload"
	"github.com/vmihailenco/msgpack"
)

func main() {
	parseCliFlags()
	configuration := parseConfiguration()

	providers, notifiers := bootstrap(configuration)

	var statistics map[string]*provider.Statistics

	if configuration.Global.CollectStats {
		loadDatabase(statistics)
	}

	scheduler := s.NewScheduler(
		providers,
		notifiers,
		statistics,
		configuration.Global.CollectStats,
	)

	scheduler.Schedule()

	failureStatusCode, _ := strconv.Atoi(configuration.Global.WebUI.FailureStatusCode)
	httpServer := http.NewHttpServer(
		scheduler,
		failureStatusCode,
		configuration.Global.WebUI.AuditLogLimit,
		configuration.Global.WebUI.ResponseLogLimit,
	)

	httpServer.Listen(configuration.Global.WebUI.Port)
}

func parseCliFlags() {
	configurationFile := flag.String("config", config.ConfigPath, "Configuration file")
	config.ConfigPath = *configurationFile

	databaseFile := flag.String("database", config.DbPath, "Database file")
	config.DbPath = *databaseFile
}

func bootstrap(configuration config.Config) ([]provider.Provider, map[string][]notifier.Notifier) {
	providerFactory := pFactory.NewProviderFactory()
	notifierFactory := nFactory.NewProviderFactory()

	providers := make([]provider.Provider, 0)
	notifiers := make(map[string][]notifier.Notifier)

	for key, value := range configuration.HealthChecks {
		providers = append(providers, providerFactory.CreateProvider(key, value))

		notifiers[key] = make([]notifier.Notifier, 0)

		for service, notifier := range value.Notifiers.Services {
			notifiers[key] = append(notifiers[key], notifierFactory.CreateNotifier(
				key,
				service,
				notifier,
				value.Notifiers.Threshold,
			))
		}
	}

	return providers, notifiers
}

func parseConfiguration() config.Config {
	log.Info(fmt.Sprintf("Open configuration file `%s`", config.ConfigPath))

	ext := path.Ext(config.ConfigPath)

	configuration := config.Config{}

	if ext == ".yaml" || ext == ".yml" {
		configuration.FromYaml(config.ConfigPath)
	} else if ext == ".json" {
		configuration.FromJson(config.ConfigPath)
	} else {
		panic("Configuration file extension can be only: json|yaml|yml")
	}

	log.Success("Parse configuration successfully")

	return configuration
}

func loadDatabase(statistics map[string]*provider.Statistics) {
	log.Info(fmt.Sprintf("Persistence mode open, trying gather database from `%s`", config.DbPath))

	database, err := ioutil.ReadFile(config.DbPath)

	if err == nil {
		msgpack.Unmarshal(database, &statistics)
	}

	log.Success("Database loaded")
}

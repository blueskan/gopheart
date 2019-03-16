package http

import (
	"encoding/json"
	"fmt"
	"github.com/blueskan/gopheart/log"
	"net/http"

	"github.com/blueskan/gopheart/provider"
)

type httpServer struct {
	scheduler provider.Scheduler
}

func NewHttpServer(scheduler provider.Scheduler, failureStatusCode int, auditLogLimit, responseLogLimit int) HttpServer {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		schedulerStatistics := scheduler.GetStatistics()

		isUnhealthy := false

		response := make(map[string]provider.StatisticsRepresentation)

		for key, val := range schedulerStatistics {
			if val.State == "Unhealthy" {
				isUnhealthy = true
			}

			response[key] = val.Representation(auditLogLimit, responseLogLimit)
		}

		statistics, _ := json.Marshal(response)

		if isUnhealthy {
			w.WriteHeader(failureStatusCode)
		}

		fmt.Fprintf(w, string(statistics))
	})

	return &httpServer{
		scheduler: scheduler,
	}
}

func (hs *httpServer) Listen(port string) {
	log.Success(fmt.Sprintf("Web Server listening on port %s", port))

	http.ListenAndServe(":"+port, nil)
}

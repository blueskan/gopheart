package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/blueskan/gopheart/provider"
)

type httpServer struct {
	scheduler provider.Scheduler
}

func NewHttpServer(scheduler provider.Scheduler, failureStatusCode int) HttpServer {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		schedulerStatistics := scheduler.GetStatistics()

		isUnhealthy := false

		for _, val := range schedulerStatistics {
			if val.State == "Unhealthy" {
				isUnhealthy = true
				break
			}
		}

		statistics, _ := json.Marshal(schedulerStatistics)

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
	fmt.Println("Web Server listening on port " + port + "!")

	http.ListenAndServe(":"+port, nil)
}

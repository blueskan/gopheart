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

func NewHttpServer(scheduler provider.Scheduler) HttpServer {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		schedulerStatistics := scheduler.GetStatistics()
		statistics, _ := json.Marshal(schedulerStatistics)

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

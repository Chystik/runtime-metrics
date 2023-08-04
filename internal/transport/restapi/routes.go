package restapi

import (
	"net/http"

	"github.com/Chystik/runtime-metrics/internal/adapters"
)

func registerHandlers(router *http.ServeMux, handlers adapters.ServerHandlers) {
	router.HandleFunc("/update/", handlers.UpdateMetric)
}

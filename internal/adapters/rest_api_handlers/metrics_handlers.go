package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Chystik/runtime-metrics/internal/models"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
)

type MetricsHandlers interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request)
}

type metricsHandlers struct {
	metricsService metricsservice.MetricsService
}

func NewMetricHandlers(ms metricsservice.MetricsService) MetricsHandlers {
	return &metricsHandlers{metricsService: ms}
}

func (mh *metricsHandlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var (
		pathRaw string
		path    []string
		metric  models.Metric
		err     error
	)

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pathRaw = strings.TrimLeft(r.URL.Path, "/update")
	path = strings.Split(pathRaw, "/")

	if len(path) != 3 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	metric.Name = path[1]

	switch path[0] {
	case "gauge":
		var val float64

		val, err = strconv.ParseFloat(path[2], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metric.Gauge = models.Gauge(val)
		mh.metricsService.UpdateGauge(metric)
	case "counter":
		var val int

		val, err = strconv.Atoi(path[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metric.Counter = models.Counter(val)
		mh.metricsService.UpdateCounter(metric)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

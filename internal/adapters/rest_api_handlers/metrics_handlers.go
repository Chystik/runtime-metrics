package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/Chystik/runtime-metrics/internal/models"
	metricsservice "github.com/Chystik/runtime-metrics/internal/service/server"
)

const tplStr = `<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Gauge Value</th>
            <th>Counter Value</th>
        </tr>
    </thead>
    <tbody>
        {{range . }}
            <tr>
                <td>{{ .Name }}</td>
                <td>{{ .Gauge }}</td>
                <td>{{ .Counter }}</td>
            </tr>
        {{ end }}
    </tbody>
</table>`

type MetricsHandlers interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request)
	GetMetric(w http.ResponseWriter, r *http.Request)
	AllMetrics(w http.ResponseWriter, r *http.Request)
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

func (mh *metricsHandlers) GetMetric(w http.ResponseWriter, r *http.Request) {
	var (
		pathRaw    string
		path       []string
		metricName string
		metric     models.Metric
		result     string
		err        error
	)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pathRaw = strings.TrimLeft(r.URL.Path, "/value")
	path = strings.Split(pathRaw, "/")

	if len(path) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	metricName = path[1]

	switch path[0] {
	case "gauge":
		metric, err = mh.metricsService.GetMetric(metricName)
		result = fmt.Sprintf("%f", metric.Gauge)
	case "counter":
		metric, err = mh.metricsService.GetMetric(metricName)
		result = fmt.Sprintf("%f", float64(metric.Counter))
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(result))
	if err != nil {
		panic(err)
	}
}

func (sh *metricsHandlers) AllMetrics(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.New("table").Parse(tplStr)
	if err != nil {
		panic(err)
	}

	allMetrics := sh.metricsService.GetAllMetrics()
	sort.Slice(allMetrics, func(i, j int) bool {
		return allMetrics[i].Name < allMetrics[j].Name
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tpl.Execute(w, allMetrics)
	if err != nil {
		panic(err)
	}
}

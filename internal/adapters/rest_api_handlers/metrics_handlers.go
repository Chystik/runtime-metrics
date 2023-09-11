package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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
            <th>Type</th>
            <th>Value</th>
        </tr>
    </thead>
    <tbody>
        {{range . }}
            <tr>
                <td>{{ .Name }}</td>
                <td>{{ .Type }}</td>
                <td>{{ .Value }}</td>
            </tr>
        {{ end }}
    </tbody>
</table>`

type metricsHandlers struct {
	metricsService metricsservice.MetricsService
}

func NewMetricsHandlers(ms metricsservice.MetricsService) *metricsHandlers {
	h := &metricsHandlers{metricsService: ms}
	return h
}

func (mh *metricsHandlers) UpdateMetric(w http.ResponseWriter, r *http.Request) {
	var (
		pathRaw string
		path    []string
		metric  models.Metric
		err     error
	)
	ctx := context.Background()

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
	metric.ID = path[1]

	switch path[0] {
	case "gauge":
		v := new(float64)

		*v, err = strconv.ParseFloat(path[2], 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		metric.Value = v
		err = mh.metricsService.UpdateGauge(ctx, metric)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "counter":
		var val int
		var v = new(int64)

		val, err = strconv.Atoi(path[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		*v = int64(val)

		metric.Delta = v
		err = mh.metricsService.UpdateCounter(ctx, metric)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
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
		metricType string
		metric     models.Metric
		result     string
		err        error
	)
	ctx := context.Background()

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
	metricType = path[0]

	metric, err = mh.metricsService.GetMetric(ctx, models.Metric{ID: metricName, MType: metricType})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	switch metricType {
	case "gauge":
		result = strconv.FormatFloat(*metric.Value, 'f', -1, 64)
	case "counter":
		result = strconv.FormatInt(*metric.Delta, 10)
	default:
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(result))
	if err != nil {
		log.Println(err)
	}
}

func (mh *metricsHandlers) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		buf    bytes.Buffer
		err    error
	)
	ctx := context.Background()

	err = json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case "gauge":
		err = mh.metricsService.UpdateGauge(ctx, metric)
	case "counter":
		err = mh.metricsService.UpdateCounter(ctx, metric)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(&buf).Encode(metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
	}
}

func (mh *metricsHandlers) GetMetricJSON(w http.ResponseWriter, r *http.Request) {
	var (
		metric models.Metric
		buf    bytes.Buffer
		err    error
	)
	ctx := context.Background()

	err = json.NewDecoder(r.Body).Decode(&metric)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := mh.metricsService.GetMetric(ctx, metric)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = json.NewEncoder(&buf).Encode(m)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf.Bytes())
	if err != nil {
		log.Println(err)
	}
}

func (mh *metricsHandlers) AllMetrics(w http.ResponseWriter, r *http.Request) {
	type formatMetrics struct {
		Name, Type, Value string
	}
	ctx := context.Background()

	var fm []formatMetrics

	tpl, err := template.New("table").Parse(tplStr)
	if err != nil {
		log.Println(err)
	}

	m, err := mh.metricsService.GetAllMetrics(ctx)
	if err != nil {
		log.Println(err)
	}

	for i := range m {
		var v string

		if m[i].Value != nil {
			v = strconv.FormatFloat(*m[i].Value, 'f', -1, 64)
		}
		if m[i].Delta != nil {
			v = strconv.FormatInt(*m[i].Delta, 10)
		}

		fm = append(fm, formatMetrics{Name: m[i].ID, Type: m[i].MType, Value: v})
	}

	sort.Slice(fm, func(i, j int) bool {
		return fm[i].Name < fm[j].Name
	})

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	err = tpl.Execute(w, fm)
	if err != nil {
		log.Println(err)
	}
}

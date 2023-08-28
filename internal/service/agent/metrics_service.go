package agentservice

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	"github.com/Chystik/runtime-metrics/internal/models"
)

type AgentService interface {
	UpdateMetrics()
	ReportMetrics()
}

type agentService struct {
	collectableMetrics config.CollectableMetrics
	runtimeMetrics     runtime.MemStats
	cache              map[string]models.Metric
	client             adapters.AgentHTTPClient
}

func New(c adapters.AgentHTTPClient, cm config.CollectableMetrics) *agentService {
	cache := make(map[string]models.Metric)

	cache["PollCount"] = models.Metric{Id: "PollCount", MType: "counter", Delta: new(int64)}
	cache["RandomValue"] = models.Metric{Id: "RandomValue", MType: "gauge", Value: new(float64)}
	for i := range cm {
		cache[cm[i]] = models.Metric{Id: cm[i], Value: new(float64)}
	}

	return &agentService{
		collectableMetrics: cm,
		runtimeMetrics:     runtime.MemStats{},
		cache:              cache,
		client:             c,
	}
}

func (as *agentService) UpdateMetrics() {
	runtime.ReadMemStats(&as.runtimeMetrics)

	for i := range as.collectableMetrics {
		r := reflect.ValueOf(as.runtimeMetrics)
		f := r.FieldByName(as.collectableMetrics[i])
		v := f.Interface()

		m, ok := as.cache[as.collectableMetrics[i]]
		if !ok {
			continue
		}

		switch val := v.(type) {
		case float64:
			m.MType = "gauge"
			*m.Value = val
		case uint64:
			m.MType = "gauge"
			*m.Value = float64(val)
		case uint32:
			m.MType = "gauge"
			*m.Value = float64(val)
		}
		m.Id = as.collectableMetrics[i]
		as.cache[as.collectableMetrics[i]] = m
	}

	pc, ok := as.cache["PollCount"]
	if ok {
		*pc.Delta += 1
		as.cache["PollCount"] = pc
	}

	rv, ok := as.cache["RandomValue"]
	if ok {
		*rv.Value = float64(rand.Intn(1000))
		as.cache["RandomValue"] = rv
	}

}

func (as *agentService) ReportMetrics() {
	err := as.client.ReportMetricsJSON(as.cache)
	if err != nil {
		fmt.Println(err)
	}
}

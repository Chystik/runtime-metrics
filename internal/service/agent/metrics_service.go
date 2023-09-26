package agentservice

import (
	"context"
	"math/rand"
	"reflect"
	"runtime"
	"sync"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/adapters"
	"github.com/Chystik/runtime-metrics/internal/models"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type AgentService interface {
	UpdateMetrics()
	UpdateGoPsUtilMetrics()
	ReportMetrics(context.Context) error
}

type agentService struct {
	collectableMetrics config.CollectableMetrics
	runtimeMetrics     runtime.MemStats
	memMetrics         *mem.VirtualMemoryStat
	cpuMetrics         []cpu.InfoStat
	mu                 sync.RWMutex
	cache              map[string]models.Metric
	client             adapters.AgentHTTPClient
}

func New(c adapters.AgentHTTPClient, cm config.CollectableMetrics) *agentService {
	cache := make(map[string]models.Metric)

	cache["PollCount"] = models.Metric{ID: "PollCount", MType: "counter", Delta: new(int64)}
	cache["RandomValue"] = models.Metric{ID: "RandomValue", MType: "gauge", Value: new(float64)}
	cache["TotalMemory"] = models.Metric{ID: "TotalMemory", Value: new(float64)}
	cache["FreeMemory"] = models.Metric{ID: "FreeMemory", Value: new(float64)}
	cache["CPUutilization1"] = models.Metric{ID: "CPUutilization1", Value: new(float64)}
	for i := range cm {
		cache[cm[i]] = models.Metric{ID: cm[i], Value: new(float64)}
	}

	return &agentService{
		collectableMetrics: cm,
		runtimeMetrics:     runtime.MemStats{},
		memMetrics:         &mem.VirtualMemoryStat{},
		cpuMetrics:         []cpu.InfoStat{},
		cache:              cache,
		client:             c,
	}
}

func (as *agentService) UpdateMetrics() {
	as.mu.Lock()
	defer as.mu.Unlock()

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
		m.ID = as.collectableMetrics[i]
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

func (as *agentService) UpdateGoPsUtilMetrics() {
	as.mu.Lock()
	defer as.mu.Unlock()

	as.memMetrics, _ = mem.VirtualMemory()
	as.cpuMetrics, _ = cpu.Info()

	totalMemory := float64(as.memMetrics.Total)
	as.cache["TotalMemory"] = models.Metric{ID: "TotalMemory", MType: "gauge", Value: &totalMemory}
	freeMemory := float64(as.memMetrics.Free)
	as.cache["FreeMemory"] = models.Metric{ID: "FreeMemory", MType: "gauge", Value: &freeMemory}

	cpuUtil1 := float64(len(as.cpuMetrics))
	as.cache["CPUutilization1"] = models.Metric{ID: "CPUutilization1", MType: "gauge", Value: &cpuUtil1}
}

func (as *agentService) ReportMetrics(ctx context.Context) error {
	as.mu.RLock()
	defer as.mu.RUnlock()

	return as.client.ReportMetricsJSONBatch(ctx, as.cache)
}

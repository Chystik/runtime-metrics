package agentservice

import (
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
	cache              map[string]interface{}
	client             adapters.AgentClient
}

func New(c adapters.AgentClient, cm config.CollectableMetrics) AgentService {
	cache := make(map[string]interface{})
	cache["PollCount"] = models.Counter(0)

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

		var val interface{}

		switch v := v.(type) {
		case float64:
			val = models.Gauge(v)
		case uint64:
			val = models.Gauge(v)
		case uint32:
			val = models.Gauge(v)
		}

		as.cache[as.collectableMetrics[i]] = val
	}

	v, ok := as.cache["PollCount"]
	if !ok {
		panic("can't get PollCount from cache")
	}

	as.cache["PollCount"] = models.Counter(v.(models.Counter) + 1)
	as.cache["RandomValue"] = models.Gauge(rand.Intn(1000))

	/* for k, v := range as.cache {
		fmt.Printf("%s \t %T: %v\n", k, v, v)
	}
	fmt.Println(as.cache["PollCount"], "================================") */
}

func (as *agentService) ReportMetrics() {
	err := as.client.ReportMetrics(as.cache)
	if err != nil {
		panic(err)
	}
}

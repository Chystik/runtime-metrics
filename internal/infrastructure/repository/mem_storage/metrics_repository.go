package memstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
)

var ErrNotFoundMetric = errors.New("not found in repository")

type memStorage struct {
	data    map[string]models.Metric
	mu      sync.Mutex
	file    *os.File
	encoder *json.Encoder
	decoder *json.Decoder
	ticker  *time.Ticker
	quit    chan bool
}

func New(cfg config.ServerConfig) (*memStorage, error) {
	var err error
	ms := &memStorage{}
	ms.data = make(map[string]models.Metric)

	if cfg.FileStoragePath != "" {
		ms.file, err = os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return nil, err
		}

		ms.ticker = time.NewTicker(time.Duration(cfg.StoreInterval))

		ms.encoder = json.NewEncoder(ms.file)
		ms.decoder = json.NewDecoder(ms.file)

		if cfg.Restore {
			err = ms.readData()
			if err != nil && err.Error() != "EOF" {
				return nil, err
			}
		}
	}

	ms.quit = make(chan bool)

	return ms, nil
}

func (ms *memStorage) UpdateGauge(metric models.Metric) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m, ok := ms.data[metric.ID]
	if !ok {
		ms.data[metric.ID] = metric
	} else {
		m.Value = metric.Value
		ms.data[metric.ID] = m
	}
}

func (ms *memStorage) UpdateCounter(metric models.Metric) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m, ok := ms.data[metric.ID]
	if !ok {
		ms.data[metric.ID] = metric
	} else {
		m.Delta = metric.Delta
		ms.data[metric.ID] = m
	}
}

func (ms *memStorage) Get(ID string) (models.Metric, error) {
	m, ok := ms.data[ID]
	if !ok {
		return models.Metric{ID: ID, MType: "", Delta: nil, Value: nil}, fmt.Errorf("metric with ID %s %w", ID, ErrNotFoundMetric)
	}

	return m, nil
}

func (ms *memStorage) GetAll() []models.Metric {
	var metrics []models.Metric

	for _, v := range ms.data {
		m := v
		metrics = append(metrics, m)
	}

	return metrics
}

func (ms *memStorage) Shutdown() error {
	if ms.file == nil {
		return nil
	}

	ms.quit <- true

	defer ms.file.Close()
	return ms.writeData()
}

func (ms *memStorage) writeData() error {
	err := ms.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = ms.file.Seek(0, 0)
	if err != nil {
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	err = ms.encoder.Encode(ms.data)
	if err != nil {
		return err
	}

	return nil
}

func (ms *memStorage) readData() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	err := ms.decoder.Decode(&ms.data)
	if err != nil {
		return err
	}

	return nil
}

func (ms *memStorage) SyncData() error {
	for {
		select {
		case <-ms.ticker.C:
			err := ms.writeData()
			if err != nil {
				return err
			}
		case <-ms.quit:
			ms.ticker.Stop()
			return nil
		}
	}
}

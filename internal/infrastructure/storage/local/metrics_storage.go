package localfs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service"
)

type localStorage struct {
	metricsRepo service.MetricsRepository
	file        *os.File
	encoder     *json.Encoder
	decoder     *json.Decoder
}

func NewMetricsStorage(cfg *config.ServerConfig, repo service.MetricsRepository) (*localStorage, error) {
	if cfg.FileStoragePath == "" {
		return nil, fmt.Errorf("file path not specified in server config: %v", cfg)
	}

	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	encoder := json.NewEncoder(file)
	decoder := json.NewDecoder(file)

	return &localStorage{
		metricsRepo: repo,
		file:        file,
		encoder:     encoder,
		decoder:     decoder,
	}, nil
}

func (ls *localStorage) Read() error {
	var m []models.Metric

	err := ls.decoder.Decode(&m)
	if err != nil {
		return err
	}

	return ls.metricsRepo.UpdateAll(context.Background(), m)
}

func (ls *localStorage) Write() error {
	err := ls.file.Truncate(0)
	if err != nil {
		return err
	}

	_, err = ls.file.Seek(0, 0)
	if err != nil {
		return err
	}

	m, err := ls.metricsRepo.GetAll(context.Background())
	if err != nil {
		return err
	}

	return ls.encoder.Encode(m)
}

func (ls *localStorage) CloseFile() error {
	if ls.file == nil {
		return nil
	}
	return ls.file.Close()
}

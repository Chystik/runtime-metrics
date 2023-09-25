package localfs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Chystik/runtime-metrics/config"
	memstorage "github.com/Chystik/runtime-metrics/internal/infrastructure/repository/mem_storage"
)

type localStorage struct {
	inMemRepo *memstorage.MemStorage
	file      *os.File
	encoder   *json.Encoder
	decoder   *json.Decoder
}

func New(cfg *config.ServerConfig, inMemRepo *memstorage.MemStorage) (*localStorage, error) {
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
		inMemRepo: inMemRepo,
		file:      file,
		encoder:   encoder,
		decoder:   decoder,
	}, nil
}

func (ls *localStorage) Read() error {
	ls.inMemRepo.Mu.RLock()
	defer ls.inMemRepo.Mu.RUnlock()

	return ls.decoder.Decode(&ls.inMemRepo.Data)
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

	ls.inMemRepo.Mu.Lock()
	defer ls.inMemRepo.Mu.Unlock()
	return ls.encoder.Encode(ls.inMemRepo.Data)
}

func (ls *localStorage) CloseFile() error {
	if ls.file == nil {
		return nil
	}
	return ls.file.Close()
}

package localfs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service"
	"github.com/Chystik/runtime-metrics/internal/service/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	errTruncate = errors.New("truncate error")
	errSeek     = errors.New("seek error")
)

type fileMock struct {
	bytes.Buffer
	io.ReaderAt
	io.Seeker
	os.FileInfo
}

func (fm fileMock) Truncate(size int64) error {
	fm.Buffer.Truncate(int(size))
	return nil
}

func (fm fileMock) Stat() (os.FileInfo, error) {
	return fm.FileInfo, nil
}

func (fm fileMock) Close() error {
	return nil
}

func (fm fileMock) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

type fileMockTruncateErr struct {
	fileMock
}

func (fm fileMockTruncateErr) Truncate(size int64) error {
	return errTruncate
}

type fileMockSeekErr struct {
	fileMock
}

func (fm fileMockSeekErr) Seek(offset int64, whence int) (int64, error) {
	return 0, errSeek
}

func newFsMock(f file) (*localStorage, *mocks.MetricsRepository) {
	repo := &mocks.MetricsRepository{}
	storage := newMetricsStorageMock(f, repo)

	return storage, repo
}

func newMetricsStorageMock(f file, repo service.MetricsRepository) *localStorage {
	encoder := json.NewEncoder(f)
	decoder := json.NewDecoder(f)

	return &localStorage{
		metricsRepo: repo,
		file:        f,
		encoder:     encoder,
		decoder:     decoder,
	}
}

func Test_CreateMetricsStorage_WithEmptyFilePath(t *testing.T) {
	t.Parallel()

	_, err := NewMetricsStorage(&config.ServerConfig{}, &mocks.MetricsRepository{})
	assert.Error(t, err)
}

func Test_CreateMetricsStorage_WithWrongFilePath(t *testing.T) {
	t.Parallel()

	cfg := &config.ServerConfig{
		FileStoragePath: "/",
	}

	_, err := NewMetricsStorage(cfg, &mocks.MetricsRepository{})
	assert.Error(t, err)
}

func Test_CreateMetricsStorage(t *testing.T) {
	t.Parallel()

	tmpFile, errCreateTmp := os.CreateTemp("", "")
	if errCreateTmp != nil {
		t.Error(errCreateTmp)
	}
	defer os.Remove(tmpFile.Name())

	cfg := &config.ServerConfig{
		FileStoragePath: tmpFile.Name(),
	}

	_, err := NewMetricsStorage(cfg, &mocks.MetricsRepository{})
	assert.NoError(t, err)
}

func Test_Read_WhenDecoderReturnsError(t *testing.T) {
	t.Parallel()

	mockStorage := newMetricsStorageMock(&fileMock{}, &mocks.MetricsRepository{})
	err := mockStorage.Read()

	assert.Error(t, err)
}

func Test_Read_WhenRepoReturnsResult(t *testing.T) {
	t.Parallel()

	mockStorage, mockRepo := newFsMock(&fileMock{})

	_, errRead := mockStorage.file.Write(encodeMetrics(generateMetrics(10)))
	assert.NoError(t, errRead)

	mockRepo.EXPECT().UpdateList(mock.Anything, mock.Anything).Return(nil)
	err := mockStorage.Read()

	assert.NoError(t, err)
}

func Test_Write_WhenTruncateReturnsError(t *testing.T) {
	t.Parallel()

	mockStorage, _ := newFsMock(&fileMockTruncateErr{})
	err := mockStorage.Write()

	assert.Error(t, err)
}

func Test_Write_WhenSeekReturnsError(t *testing.T) {
	t.Parallel()

	mockStorage, _ := newFsMock(&fileMockSeekErr{})
	err := mockStorage.Write()

	assert.Error(t, err)
}

func Test_Write_WhenRepoReturnsError(t *testing.T) {
	t.Parallel()

	mockStorage, mockRepo := newFsMock(&fileMock{})

	mockRepo.EXPECT().GetAll(mock.Anything).Return([]models.Metric{}, errors.New("some repo err"))
	err := mockStorage.Write()

	assert.Error(t, err)
}

func Test_Write_WhenRepoReturnsResult(t *testing.T) {
	t.Parallel()

	mockStorage, mockRepo := newFsMock(&fileMock{})

	mockRepo.EXPECT().GetAll(mock.Anything).Return([]models.Metric{}, nil)
	err := mockStorage.Write()

	assert.NoError(t, err)
}

func Test_CloseFile_WhenFileIsNil(t *testing.T) {
	t.Parallel()

	mockStorage, _ := newFsMock(nil)
	err := mockStorage.CloseFile()

	assert.Error(t, err)
}

func Test_CloseFile(t *testing.T) {
	t.Parallel()

	mockStorage, _ := newFsMock(&fileMock{})
	err := mockStorage.CloseFile()

	assert.NoError(t, err)
}

func generateMetrics(count int) []models.Metric {
	m := make([]models.Metric, count)

	rand.New(rand.NewSource(time.Now().UnixNano()))

	randMetricType := [2]string{"gauge", "counter"}

	for i := range m {
		mName := fmt.Sprintf("TestMetric%d", i)
		mType := randMetricType[rand.Intn(2)]
		m[i] = generateMetric(mType, mName)
	}

	return m
}

func generateMetric(metricType string, metricName string) models.Metric {
	var m models.Metric

	min := 1e1
	max := 1e3

	m.ID = metricName
	m.MType = metricType

	switch metricType {
	case "gauge":
		m.Value = new(float64)
		*m.Value = min + rand.Float64()*(max-min)
	case "counter":
		m.Delta = new(int64)
		*m.Delta = int64(min + rand.Float64()*(max-min))
	}

	return m
}

func encodeMetrics(metrics []models.Metric) []byte {
	var m bytes.Buffer

	err := json.NewEncoder(&m).Encode(metrics)
	if err != nil {
		log.Fatal(err)
	}

	return m.Bytes()
}

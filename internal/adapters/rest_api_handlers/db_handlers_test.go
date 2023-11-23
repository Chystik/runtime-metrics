package handlers

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Chystik/runtime-metrics/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const pingTimeOut = 3 * time.Second

func TestNewDBHandlers(t *testing.T) {
	t.Parallel()

	handlers, _ := getDBHandlersMocks()

	assert.NotNil(t, handlers)
}

type dbHandlersMocks struct {
	db          *mocks.PostgresClient
	logger      *mocks.AppLogger
	pingTimeout time.Duration
}

func getDBHandlersMocks() (*dbHandlers, *dbHandlersMocks) {
	m := &dbHandlersMocks{
		db:          &mocks.PostgresClient{},
		logger:      &mocks.AppLogger{},
		pingTimeout: pingTimeOut,
	}
	handlers := NewDBHandlers(m.db, m.logger, pingTimeOut)

	return handlers, m
}

func Test_dbHandlers_PingDB(t *testing.T) {
	t.Parallel()

	handlers, mks := getDBHandlersMocks()

	expStatus := http.StatusOK
	expContextType := "text/plain"

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	mks.db.EXPECT().Ping(mock.Anything).Return(nil)
	handlers.PingDB(rec, req)
	res := rec.Result()

	assert.Equal(t, expStatus, res.StatusCode)

	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.Equal(t, expContextType, res.Header.Get("Content-Type"))
}

func Test_dbHandlers_PingDB_DBReturnsError(t *testing.T) {
	t.Parallel()

	handlers, mks := getDBHandlersMocks()

	expStatus := http.StatusInternalServerError
	expContextType := "text/plain"
	expError := errors.New("some error")

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	mks.db.EXPECT().Ping(mock.Anything).Return(expError)
	mks.logger.EXPECT().Error(expError.Error())

	handlers.PingDB(rec, req)
	res := rec.Result()

	assert.Equal(t, expStatus, res.StatusCode)

	defer res.Body.Close()
	_, err := io.ReadAll(res.Body)

	assert.NoError(t, err)
	assert.Equal(t, expContextType, res.Header.Get("Content-Type"))
}

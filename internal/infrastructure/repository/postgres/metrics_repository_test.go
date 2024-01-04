package postgresrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service"
	"github.com/Chystik/runtime-metrics/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type conRetryer struct{}

func (cr conRetryer) DoWithRetry(fn func() error) error {
	return fn()
}

func newConRetryer() service.ConnectionRetrier {
	return conRetryer{}
}

func newSqlxDB(t *testing.T) (*sqlx.DB, sqlmock.Sqlmock) {
	mockDB, sqlMock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
		return nil, nil
	}
	return sqlx.NewDb(mockDB, "sqlmock"), sqlMock
}

func Test_UpdateGauge(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	m := models.Metric{ID: "test", MType: "gauge", Value: new(float64)}
	query := regexp.QuoteMeta(`
			INSERT INTO	praktikum.metrics (id, m_type, m_value)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET 
				m_value = EXCLUDED.m_value`)

	mockSQL.ExpectExec(query).WithArgs(m.ID, m.MType, m.Value).WillReturnResult(sqlmock.NewResult(0, 0))
	err := conRetMock.DoWithRetry(func() error {
		return nil
	})
	assert.NoError(t, err)
	logMock.EXPECT().Error(mock.Anything).Return()

	err = pgRepo.UpdateGauge(context.Background(), m)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_UpdateGauge_WhenRetryerReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, _ := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	err := conRetMock.DoWithRetry(func() error {
		return errors.New("some err")
	})
	assert.Error(t, err)

	logMock.EXPECT().Error(mock.Anything).Return()

	err = pgRepo.UpdateGauge(context.Background(), models.Metric{})

	assert.Error(t, err)
}

func Test_UpdateCounter(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	m := models.Metric{ID: "test", MType: "gauge", Delta: new(int64)}
	query := regexp.QuoteMeta(`
			INSERT INTO	praktikum.metrics (id, m_type, m_delta)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO
			UPDATE SET
				m_delta = $3 + (SELECT m_delta
					FROM praktikum.metrics
					WHERE id = $1)`)

	mockSQL.ExpectExec(query).WithArgs(m.ID, m.MType, m.Delta).WillReturnResult(sqlmock.NewResult(1, 1))
	err := conRetMock.DoWithRetry(func() error {
		return nil
	})
	assert.NoError(t, err)
	logMock.EXPECT().Error(mock.Anything).Return()

	err = pgRepo.UpdateCounter(context.Background(), m)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_UpdateCounter_WhenRetryerReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, _ := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	err := conRetMock.DoWithRetry(func() error {
		return errors.New("some err")
	})
	assert.Error(t, err)

	logMock.EXPECT().Error(mock.Anything).Return()

	err = pgRepo.UpdateCounter(context.Background(), models.Metric{})

	assert.Error(t, err)
}

func Test_Get(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	m := generateMetric("test", "counter")
	query := regexp.QuoteMeta(`
			SELECT id, m_type, m_value, m_delta
			FROM praktikum.metrics
			WHERE id = $1`)

	err := conRetMock.DoWithRetry(func() error {
		var v, d string

		if m.Value != nil {
			v = strconv.Itoa(int(*m.Value))
		}
		if m.Delta != nil {
			d = strconv.Itoa(int(*m.Delta))
		}

		mockSQL.ExpectQuery(query).WithArgs(m.ID).WillReturnRows(
			sqlmock.NewRows([]string{
				"id",
				"m_type",
				"m_value",
				"m_delta",
			}).AddRow(
				m.ID,
				m.MType,
				v,
				d,
			))

		return nil
	})
	assert.NoError(t, err)
	logMock.EXPECT().Error(mock.Anything).Return()

	_, err = pgRepo.Get(context.Background(), m)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_Get_WhenRetryerReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, _ := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := &mocks.ConnectionRetrier{}
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	conRetMock.EXPECT().DoWithRetry(mock.Anything).Return(sql.ErrNoRows)
	logMock.EXPECT().Error(mock.Anything).Return()

	_, err := pgRepo.Get(context.Background(), models.Metric{})

	assert.Error(t, err)
}

func Test_GetAll(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	query := regexp.QuoteMeta(`
			SELECT id, m_type, m_value, m_delta
			FROM praktikum.metrics`)

	mockSQL.ExpectQuery(query).WithArgs().WillReturnRows(sqlmock.NewRows([]string{"id", "m_type"}).AddRow(1, 1))
	err := conRetMock.DoWithRetry(func() error {
		return nil
	})
	assert.NoError(t, err)
	logMock.EXPECT().Error(mock.Anything).Return()

	_, err = pgRepo.GetAll(context.Background())

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_GetAll_WhenRetryerReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, _ := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := &mocks.ConnectionRetrier{}
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	conRetMock.EXPECT().DoWithRetry(mock.Anything).Return(sql.ErrNoRows)
	logMock.EXPECT().Error(mock.Anything).Return()

	_, err := pgRepo.GetAll(context.Background())

	assert.Error(t, err)
}

func Test_UpdateList(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)

	metrics := generateMetrics(10)
	query := regexp.QuoteMeta(`
		INSERT INTO	praktikum.metrics (id, m_type, m_value, m_delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO
		UPDATE SET
			m_value = EXCLUDED.m_value,
			m_delta = $4 + (SELECT m_delta
			FROM praktikum.metrics
			WHERE id = $1)`)

	err := conRetMock.DoWithRetry(func() error {
		mockSQL.ExpectBegin()
		return nil
	})
	assert.NoError(t, err)

	stmt := mockSQL.ExpectPrepare(query)
	for _, m := range metrics {
		stmt.ExpectExec().WithArgs(
			m.ID,
			m.MType,
			m.Value,
			m.Delta,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mockSQL.ExpectCommit()
	logMock.EXPECT().Error(mock.Anything).Return()

	err = pgRepo.UpdateList(context.Background(), metrics)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.NoError(t, err)
}

func Test_UpdateList_WhenBeginReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)
	expErr := errors.New("tx begin error")

	err := conRetMock.DoWithRetry(func() error {
		mockSQL.ExpectBegin().WillReturnError(expErr)
		return expErr
	})
	assert.Error(t, err)
	logMock.EXPECT().Error(mock.Anything).Return()

	err = pgRepo.UpdateList(context.Background(), []models.Metric{})

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func Test_UpdateList_WhenPrepareContextReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)
	expErr := errors.New("prepare context error")

	metrics := generateMetrics(10)
	query := regexp.QuoteMeta(`
		INSERT INTO	praktikum.metrics (id, m_type, m_value, m_delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO
		UPDATE SET
			m_value = EXCLUDED.m_value,
			m_delta = $4 + (SELECT m_delta
			FROM praktikum.metrics
			WHERE id = $1)`)

	err := conRetMock.DoWithRetry(func() error {
		mockSQL.ExpectBegin()
		return nil
	})
	assert.NoError(t, err)

	mockSQL.ExpectPrepare(query).WillReturnError(expErr)
	logMock.EXPECT().Error(mock.Anything).Return()
	mockSQL.ExpectRollback()

	err = pgRepo.UpdateList(context.Background(), metrics)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func Test_UpdateList_WhenExecContextReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)
	expErr := errors.New("exec context error")

	metrics := generateMetrics(10)
	query := regexp.QuoteMeta(`
		INSERT INTO	praktikum.metrics (id, m_type, m_value, m_delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO
		UPDATE SET
			m_value = EXCLUDED.m_value,
			m_delta = $4 + (SELECT m_delta
			FROM praktikum.metrics
			WHERE id = $1)`)

	err := conRetMock.DoWithRetry(func() error {
		mockSQL.ExpectBegin()
		return nil
	})
	assert.NoError(t, err)

	mockSQL.ExpectPrepare(query)
	mockSQL.ExpectExec(query).WithArgs(
		metrics[0].ID,
		metrics[0].MType,
		metrics[0].Value,
		metrics[0].Delta,
	).WillReturnError(expErr)

	logMock.EXPECT().Error(mock.Anything).Return()
	mockSQL.ExpectRollback()

	err = pgRepo.UpdateList(context.Background(), metrics)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func Test_UpdateList_WhenCommitReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)
	expErr := errors.New("commit error")

	metrics := generateMetrics(10)
	query := regexp.QuoteMeta(`
		INSERT INTO	praktikum.metrics (id, m_type, m_value, m_delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO
		UPDATE SET
			m_value = EXCLUDED.m_value,
			m_delta = $4 + (SELECT m_delta
			FROM praktikum.metrics
			WHERE id = $1)`)

	err := conRetMock.DoWithRetry(func() error {
		mockSQL.ExpectBegin()
		return nil
	})
	assert.NoError(t, err)

	stmt := mockSQL.ExpectPrepare(query)
	for _, m := range metrics {
		stmt.ExpectExec().WithArgs(
			m.ID,
			m.MType,
			m.Value,
			m.Delta,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mockSQL.ExpectCommit().WillReturnError(expErr)
	logMock.EXPECT().Error(mock.Anything).Return()

	err = pgRepo.UpdateList(context.Background(), metrics)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
}

func Test_UpdateList_WhenRollBackReturnsError(t *testing.T) {
	t.Parallel()

	sqlxDB, mockSQL := newSqlxDB(t)
	defer sqlxDB.Close()

	conRetMock := newConRetryer()
	logMock := &mocks.AppLogger{}

	pgRepo := NewMetricsRepo(sqlxDB, conRetMock, logMock)
	expErr := errors.New("rollback error")

	metrics := generateMetrics(10)
	query := regexp.QuoteMeta(`
		INSERT INTO	praktikum.metrics (id, m_type, m_value, m_delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO
		UPDATE SET
			m_value = EXCLUDED.m_value,
			m_delta = $4 + (SELECT m_delta
			FROM praktikum.metrics
			WHERE id = $1)`)

	err := conRetMock.DoWithRetry(func() error {
		mockSQL.ExpectBegin()
		return nil
	})
	assert.NoError(t, err)

	mockSQL.ExpectPrepare(query).WillReturnError(expErr)
	logMock.EXPECT().Error(mock.Anything).Return()
	mockSQL.ExpectRollback().WillReturnError(expErr)

	err = pgRepo.UpdateList(context.Background(), metrics)

	assert.NoError(t, mockSQL.ExpectationsWereMet())
	assert.Error(t, err)
	assert.Equal(t, expErr, err)
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

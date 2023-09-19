package postgres

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"syscall"
	"time"

	"github.com/Chystik/runtime-metrics/internal/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	attempts      = 3
	triesInterval = 1
	deltaInterval = 2

	ErrNotFoundMetric = errors.New("not found in repository")

	logRetryConnection = "cannot connect to the database server. next try in %d seconds\n"
)

type pgRepo struct {
	db *sqlx.DB
	l  *zap.Logger
}

func NewMetricsRepo(db *sqlx.DB, logger *zap.Logger) *pgRepo {
	return &pgRepo{
		db: db,
		l:  logger,
	}
}

func (pg *pgRepo) UpdateGauge(ctx context.Context, metric models.Metric) error {
	var err error

	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_value)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET 
				m_value = EXCLUDED.m_value`

	doWithRetry(pg.l, func() error {
		_, err = pg.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value)
		return err
	})
	if err != nil {
		pg.l.Error(err.Error())
		return err
	}

	return nil
}

func (pg *pgRepo) UpdateCounter(ctx context.Context, metric models.Metric) error {
	var err error

	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_delta)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET 
				m_delta = $3 + (SELECT m_delta
					FROM praktikum.metrics
					WHERE id = $1)`

	doWithRetry(pg.l, func() error {
		_, err := pg.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Delta)
		return err
	})
	if err != nil {
		pg.l.Error(err.Error())
		return err
	}

	return nil
}

func (pg *pgRepo) Get(ctx context.Context, metric models.Metric) (models.Metric, error) {
	var m models.Metric
	var err error

	query := `
			SELECT id, m_type, m_value, m_delta
			FROM praktikum.metrics
			WHERE id = $1`

	doWithRetry(pg.l, func() error {
		return pg.db.GetContext(ctx, &m, query, metric.ID)
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return m, ErrNotFoundMetric
		}
	}

	return m, nil
}

func (pg *pgRepo) GetAll(ctx context.Context) ([]models.Metric, error) {
	var metrics []models.Metric
	var rows *sql.Rows
	var err error

	query := `
			SELECT id, m_type, m_value, m_delta
			FROM praktikum.metrics`

	doWithRetry(pg.l, func() error {
		rows, err = pg.db.QueryContext(ctx, query)
		return err
	})
	if err != nil {
		pg.l.Error(err.Error())
		return nil, err
	}

	for rows.Next() {
		var m models.Metric

		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			pg.l.Error(err.Error())
			return nil, err
		}

		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		pg.l.Error(err.Error())
		return nil, err
	}
	rows.Close()

	return metrics, nil
}

func (pg *pgRepo) UpdateAll(ctx context.Context, metrics []models.Metric) (err error) {
	var tx *sql.Tx

	doWithRetry(pg.l, func() error {
		tx, err = pg.db.Begin()
		return err
	})
	if err != nil {
		pg.l.Error(err.Error())
		return err
	}
	defer func() {
		e := tx.Rollback()
		if e.Error() != "sql: transaction has already been committed or rolled back" {
			err = e
		}
	}()

	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_value, m_delta)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO 
			UPDATE SET 
				m_value = EXCLUDED.m_value, 
				m_delta = $4 + (SELECT m_delta
				FROM praktikum.metrics
				WHERE id = $1)`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		pg.l.Error(err.Error())
		return err
	}

	for _, m := range metrics {
		_, err = stmt.ExecContext(ctx,
			m.ID,
			m.MType,
			m.Value,
			m.Delta,
		)
		if err != nil {
			pg.l.Error(err.Error())
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		pg.l.Error(err.Error())
		return err
	}

	return nil
}

func doWithRetry(l *zap.Logger, retryableFunc func() error) {
	err := retryableFunc()
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := attempts
			ti := triesInterval
			for a > 0 && err != nil {
				l.Sugar().Infof(logRetryConnection, ti)
				time.Sleep(time.Duration(ti) * time.Second)
				err = retryableFunc()
				a--
				ti = ti + deltaInterval
			}
		}
	}
}

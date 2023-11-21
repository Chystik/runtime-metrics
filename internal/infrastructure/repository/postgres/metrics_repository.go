package postgresrepo

import (
	"context"
	"database/sql"
	"errors"
	"sort"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/Chystik/runtime-metrics/internal/service"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFoundMetric = errors.New("not found in repository")
)

type pgRepo struct {
	db *sqlx.DB
	r  service.ConnectionRetrier
	l  service.AppLogger
}

func NewMetricsRepo(db *sqlx.DB, r service.ConnectionRetrier, logger service.AppLogger) *pgRepo {
	return &pgRepo{
		db: db,
		r:  r,
		l:  logger,
	}
}

func (pg *pgRepo) UpdateGauge(ctx context.Context, metric models.Metric) error {
	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_value)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET 
				m_value = EXCLUDED.m_value`

	err := pg.r.DoWithRetry(func() error {
		_, err := pg.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value)
		return err
	})
	if err != nil {
		pg.l.Error(err.Error())
		return err
	}

	return nil
}

func (pg *pgRepo) UpdateCounter(ctx context.Context, metric models.Metric) error {
	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_delta)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET 
				m_delta = $3 + (SELECT m_delta
					FROM praktikum.metrics
					WHERE id = $1)`

	err := pg.r.DoWithRetry(func() error {
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

	query := `
			SELECT id, m_type, m_value, m_delta
			FROM praktikum.metrics
			WHERE id = $1`

	err := pg.r.DoWithRetry(func() error {
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
	var err error

	query := `
			SELECT id, m_type, m_value, m_delta
			FROM praktikum.metrics`

	err = pg.r.DoWithRetry(func() error {
		err = pg.db.SelectContext(ctx, &metrics, query)
		return err
	})
	if err != nil {
		pg.l.Error(err.Error())
		return nil, err
	}

	return metrics, nil
}

func (pg *pgRepo) UpdateAll(ctx context.Context, metrics []models.Metric) (err error) {
	var tx *sql.Tx

	sort.Slice(metrics, func(i, j int) bool { // prevent error on concurrent update: deadlock detected (SQLSTATE 40P01)
		return metrics[i].ID < metrics[j].ID
	})

	err = pg.r.DoWithRetry(func() error {
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

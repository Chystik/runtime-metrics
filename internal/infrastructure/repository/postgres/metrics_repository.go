package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/runtime-metrics/internal/models"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	attempts      = 3
	triesInterval = 1
	deltaInterval = 2

	ErrNotFoundMetric = errors.New("not found in repository")
	errConnection     = errors.New("cannot connect to the db server. exit")

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
	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_value)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO 
			UPDATE SET 
				m_value = EXCLUDED.m_value`

	_, err := pg.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value)
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

	_, err := pg.ExecContext(ctx, query, metric.ID, metric.MType, metric.Delta)
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

	err := pg.GetContext(ctx, &m, query, metric.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			pg.l.Error(err.Error())
			return m, ErrNotFoundMetric
		}
	}

	return m, nil
}

func (pg *pgRepo) GetAll(ctx context.Context) ([]models.Metric, error) {
	var metrics []models.Metric

	query := `
			SELECT id, m_type, m_value, m_delta
			FROM praktikum.metrics`

	rows, err := pg.QueryContext(ctx, query)
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
	tx, err := pg.Begin()
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

package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Chystik/runtime-metrics/internal/models"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFoundMetric = errors.New("not found in repository")
)

type pgRepo struct {
	db *sqlx.DB
}

func NewMetricsRepo(db *sqlx.DB) *pgRepo {
	return &pgRepo{
		db: db,
	}
}

func (pg *pgRepo) UpdateGauge(ctx context.Context, metric models.Metric) error {
	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_value)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE SET m_value = EXCLUDED.m_value`

	_, err := pg.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value)
	if err != nil {
		return err
	}

	return nil
}

func (pg *pgRepo) UpdateCounter(ctx context.Context, metric models.Metric) error {
	query := `
			INSERT INTO	praktikum.metrics (id, m_type, m_delta)
			VALUES ($1, $2, $3)
			ON CONFLICT (id) DO UPDATE SET m_delta = EXCLUDED.m_delta`

	_, err := pg.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Delta)
	if err != nil {
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

	err := pg.db.GetContext(ctx, &m, query, metric.ID)
	if err != nil {
		if err == sql.ErrNoRows {
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

	rows, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var m models.Metric

		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, m)
	}

	if rows.Err() != nil {
		return nil, err
	}
	rows.Close()

	return metrics, nil
}

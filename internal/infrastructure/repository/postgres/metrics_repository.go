package postgres

import (
	"github.com/jmoiron/sqlx"
)

type pgRepo struct {
	db *sqlx.DB
}

func NewMetricsRepo(db *sqlx.DB) *pgRepo {
	return &pgRepo{
		db: db,
	}
}

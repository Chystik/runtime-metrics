package db

import (
	"context"
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type pgClient struct {
	db  *sqlx.DB
	dsn string
}

func NewPgClient(cfg *config.ServerConfig) (*pgClient, error) {
	db, err := sqlx.Open("pgx", cfg.DbDsn)
	if err != nil {
		return nil, err
	}

	return &pgClient{
		db:  db,
		dsn: cfg.DbDsn,
	}, nil
}

// Connect to a database and verify with a ping
func (pc *pgClient) Connect(ctx context.Context) error {
	var err error
	pc.db, err = sqlx.ConnectContext(ctx, "pgx", pc.dsn)
	return err
}

func (pc *pgClient) Disconnect(ctx context.Context) error {
	return pc.db.Close()
}

func (pc *pgClient) Ping(ctx context.Context) error {
	return pc.db.PingContext(ctx)
}

func (pc *pgClient) PingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	w.Header().Set("Content-Tye", "text/plain")

	err := pc.Ping(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

package db

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	connStr   = "host=%s port=%d user=%s password=%s sslmode=%s"
	connStrDB = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
)

type pgClient struct {
	db *sqlx.DB
	c  *pgx.ConnConfig
}

// opens a db and perform migrations
func NewPgClient(cfg *config.ServerConfig) (*pgClient, error) {
	cc, err := pgx.ParseURI(cfg.DBDsn)
	if err != nil {
		return nil, err
	}

	db, err := sqlx.Open("pgx", cfg.DBDsn)
	if err != nil {
		return nil, err
	}

	return &pgClient{
		db: db,
		c:  &cc,
	}, nil
}

// Connect to a database and verify with a ping, if successful - create db if not exist
func (pc *pgClient) Connect(ctx context.Context) (*sqlx.DB, error) {
	var err error
	var SSLmode string

	if pc.c.TLSConfig == nil {
		SSLmode = "disable"
	}

	pc.db, err = sqlx.ConnectContext(ctx, "pgx", fmt.Sprintf(connStr, pc.c.Host, pc.c.Port, pc.c.User, pc.c.Password, SSLmode))
	if err != nil {
		return nil, err
	}

	_, err = pc.db.Exec(fmt.Sprintf("CREATE DATABASE %s", pc.c.Database))
	if err != nil {
		pgEr, ok := err.(*pgconn.PgError)
		if !ok { // not pgconn error
			return nil, err
		} else if pgEr.Code != "42P04" { // database <db_name> already exists (SQLSTATE 42P04)
			return nil, err
		}
	}

	pc.db, err = sqlx.ConnectContext(ctx, "pgx", fmt.Sprintf(connStrDB, pc.c.Host, pc.c.Port, pc.c.User, pc.c.Password, pc.c.Database, SSLmode))
	if err != nil {
		return nil, err
	}

	return pc.db, nil
}

func (pc *pgClient) Migrate() error {
	d, err := postgres.WithInstance(pc.db.DB, &postgres.Config{})
	if err != nil {
		fmt.Println("1", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		pc.c.Database, d)
	if err != nil {
		fmt.Println("2", err)
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		return err
	}

	return nil
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

package db

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Chystik/runtime-metrics/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

var (
	connStr            = "host=%s port=%d user=%s password=%s sslmode=%s"
	connStrDB          = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	logDatabaseCreated = "database %s created"
)

type pgClient struct {
	db *sqlx.DB
	c  *pgx.ConnConfig
	l  *zap.Logger
}

// opens a db and perform migrations
func NewPgClient(cfg *config.ServerConfig, logger *zap.Logger) (*pgClient, error) {
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
		l:  logger,
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
		pc.l.Error(err.Error())
		return nil, err
	}

	_, err = pc.db.Exec(fmt.Sprintf("CREATE DATABASE %s", pc.c.Database))
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgerrcode.DuplicateDatabase != pgErr.Code {
			pc.l.Error(err.Error())
			return nil, err
		}
		pc.l.Info(err.Error())
	} else {
		pc.l.Info(fmt.Sprintf(logDatabaseCreated, pc.c.Database))
	}

	pc.db, err = sqlx.ConnectContext(ctx, "pgx", fmt.Sprintf(connStrDB, pc.c.Host, pc.c.Port, pc.c.User, pc.c.Password, pc.c.Database, SSLmode))
	if err != nil {
		pc.l.Error(err.Error())
		return nil, err
	}

	return pc.db, nil
}

func (pc *pgClient) Migrate() error {
	d, err := postgres.WithInstance(pc.db.DB, &postgres.Config{})
	if err != nil {
		pc.l.Error(err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		pc.c.Database, d)
	if err != nil {
		pc.l.Error(err.Error())
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		pc.l.Error(err.Error())
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
		pc.l.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

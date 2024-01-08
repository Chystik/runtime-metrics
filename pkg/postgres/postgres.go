package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Chystik/runtime-metrics/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const defaultPgPort uint16 = 5432

var (
	connStr            = "host=%s port=%d user=%s password=%s sslmode=%s"
	connStrDB          = "host=%s port=%d user=%s password=%s dbname=%s sslmode=%s"
	logDatabaseCreated = "database %s created"
)

type Postgres struct {
	*sqlx.DB
	connConfig *pgx.ConnConfig
	logger     logger.Logger
}

// New opens a db and perform migrations
func New(uri string, logger logger.Logger) (*Postgres, error) {
	cc, err := pgx.ParseURI(uri)
	if err != nil {
		return nil, err
	}

	if cc.Port == 0 {
		cc.Port = defaultPgPort
	}

	db, err := sqlx.Open("pgx", uri)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		DB:         db,
		connConfig: &cc,
		logger:     logger,
	}, nil
}

// Connect connects to the database and verify with a ping, if successful - create db if not exist
func (pc *Postgres) Connect(ctx context.Context) error {
	var err error
	var SSLmode string

	if pc.connConfig.TLSConfig == nil {
		SSLmode = "disable"
	}

	pc.DB, err = sqlx.ConnectContext(
		ctx,
		"pgx",
		fmt.Sprintf(
			connStr,
			pc.connConfig.Host,
			pc.connConfig.Port,
			pc.connConfig.User,
			pc.connConfig.Password,
			SSLmode,
		),
	)
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	_, err = pc.DB.Exec(fmt.Sprintf("CREATE DATABASE %s", pc.connConfig.Database))
	if err != nil {
		var pgErr *pgconn.PgError
		if !errors.As(err, &pgErr) || pgerrcode.DuplicateDatabase != pgErr.Code {
			pc.logger.Error(err.Error())
			return err
		}
		pc.logger.Info(err.Error())
	} else {
		pc.logger.Info(fmt.Sprintf(logDatabaseCreated, pc.connConfig.Database))
	}

	pc.DB, err = sqlx.ConnectContext(
		ctx,
		"pgx",
		fmt.Sprintf(
			connStrDB,
			pc.connConfig.Host,
			pc.connConfig.Port,
			pc.connConfig.User,
			pc.connConfig.Password,
			pc.connConfig.Database,
			SSLmode,
		),
	)
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	return nil
}

// Migrate applies all up migrations
func (pc *Postgres) Migrate() error {
	d, err := postgres.WithInstance(pc.DB.DB, &postgres.Config{})
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://schema",
		pc.connConfig.Database, d)
	if err != nil {
		pc.logger.Error(err.Error())
		return err
	}

	err = m.Up()
	if err != nil && err.Error() != "no change" {
		pc.logger.Error(err.Error())
		return err
	}

	return nil
}

func (pc *Postgres) Disconnect(ctx context.Context) error {
	return pc.DB.Close()
}

func (pc *Postgres) Ping(ctx context.Context) error {
	return pc.PingContext(ctx)
}

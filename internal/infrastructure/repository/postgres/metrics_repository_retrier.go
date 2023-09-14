package postgres

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"syscall"
	"time"

	"github.com/Chystik/runtime-metrics/internal/models"
)

func (pg *pgRepo) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	res, err := pg.db.ExecContext(ctx, query, args...)
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := attempts
			ti := triesInterval
			for a > 0 && err != nil {
				pg.l.Sugar().Infof(logRetryConnection, ti)
				time.Sleep(time.Duration(ti) * time.Second)
				res, err = pg.db.ExecContext(ctx, query, args...)
				a--
				ti = ti + deltaInterval
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (pg *pgRepo) Begin() (*sql.Tx, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := attempts
			ti := triesInterval
			for a > 0 && err != nil {
				pg.l.Sugar().Infof(logRetryConnection, ti)
				time.Sleep(time.Duration(ti) * time.Second)
				tx, err = pg.db.Begin()
				a--
				ti = ti + deltaInterval
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return tx, nil
}

func (pg *pgRepo) GetContext(ctx context.Context, m *models.Metric, query string, id string) error {
	err := pg.db.GetContext(ctx, m, query, id)
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := attempts
			ti := triesInterval
			for a > 0 && err != nil {
				pg.l.Sugar().Infof(logRetryConnection, ti)
				time.Sleep(time.Duration(ti) * time.Second)
				err = pg.db.GetContext(ctx, m, query, id)
				a--
				ti = ti + deltaInterval
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (pg *pgRepo) QueryContext(ctx context.Context, query string) (*sql.Rows, error) {
	r, err := pg.db.QueryContext(ctx, query)
	if err != nil {
		var netOpErr *net.OpError
		if errors.As(err, &netOpErr) && errors.Is(netOpErr, syscall.ECONNREFUSED) {
			a := attempts
			ti := triesInterval
			for a > 0 && err != nil {
				pg.l.Sugar().Infof(logRetryConnection, ti)
				time.Sleep(time.Duration(ti) * time.Second)
				r, err = pg.db.QueryContext(ctx, query)
				a--
				ti = ti + deltaInterval
			}
		}
		if err != nil {
			return nil, err
		}
	}
	return r, nil
}

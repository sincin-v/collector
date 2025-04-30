package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sincin-v/collector/internal/logger"
)

type DBClient struct {
	db             *sql.DB
	retryIntervals []time.Duration
}

func New(ctx context.Context, dns string, retryIntervals []time.Duration) (*DBClient, error) {
	database, err := sql.Open("pgx", dns)
	if err != nil {
		logger.Log.Error("[DBClient] Could no create DB connection. Error: %s", err)
		return nil, err
	}

	databaseClient := &DBClient{
		db:             database,
		retryIntervals: retryIntervals,
	}
	return databaseClient, nil
}

func (d *DBClient) Execute(ctx context.Context, query string, args ...any) error {
	var err error
	for _, retryInterval := range d.retryIntervals {
		_, err = d.db.ExecContext(ctx, query, args...)
		if err == nil {
			return nil
		} else {
			logger.Log.Errorf("[DBClient] Could not execute query %s. Error: %s", query, err)
			time.Sleep(retryInterval)
		}
	}
	return fmt.Errorf("%w", fmt.Errorf("cannot execute query %s", err))
}

func (d *DBClient) RowQuery(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	var resRow *sql.Row
	var err error
	for _, retryInterval := range d.retryIntervals {
		resRow = d.db.QueryRowContext(ctx, query, args...)

		var err = resRow.Err()

		if err == nil {
			return resRow, nil
		} else {
			pgErr, ok := err.(*pgconn.PgError)
			if !ok {
				logger.Log.Errorf("[DBClient] Cannot convert error to pgError. Error: %s", err)
				break
			}
			if pgerrcode.IsConnectionException(pgErr.Code) {
				logger.Log.Errorf("[DBClient] Could not execute query %s. Error: %s", query, err)
				time.Sleep(retryInterval)
				continue
			} else {
				logger.Log.Errorf("[DBClient] Could not execute query %s. Error: %s", query, err)
				break
			}
		}
	}
	return nil, fmt.Errorf("%w", fmt.Errorf("cannot execute query %s", err))
}

func (d *DBClient) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	var err error
	var resRows *sql.Rows
	for _, retryInterval := range d.retryIntervals {
		resRows, err = d.db.QueryContext(ctx, query, args...)
		if err == nil {
			return resRows, nil
		} else {
			pgErr, ok := err.(*pgconn.PgError)
			if !ok {
				logger.Log.Errorf("[DBClient] Cannot convert error to pgError. Error: %s", err)
				break
			}
			if pgerrcode.IsConnectionException(pgErr.Code) {
				logger.Log.Errorf("[DBClient] Could not execute query %s. Error: %s", query, err)
				time.Sleep(retryInterval)
				continue
			} else {
				logger.Log.Errorf("[DBClient] Could not execute query %s. Error: %s", query, err)
				break
			}
		}
	}

	return nil, fmt.Errorf("%w", fmt.Errorf("cannot execute query %s", err))
}

func (d *DBClient) BeginTx(ctx context.Context) (*sql.Tx, error) {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		logger.Log.Error("[DBClient] Cannot start Tx. Error: %s", err)
		return nil, err
	}
	return tx, nil
}

func (d *DBClient) Close() {
	if err := d.db.Close(); err != nil {
		logger.Log.Error("[DBClient] Could not close db connection. Error:  %s", err)
	}
}

func (d *DBClient) Ping(ctx context.Context) (bool, error) {
	for _, retryInterval := range d.retryIntervals {
		err := d.db.PingContext(ctx)
		if err != nil {
			logger.Log.Error("[DBClient] Ping. Error: %s", err)
			time.Sleep(retryInterval)
			continue
		} else {
			return true, nil
		}
	}
	return false, fmt.Errorf("%w", errors.New("cannot connect to DB. Ping is failed"))
}

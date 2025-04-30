package db

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sincin-v/collector/internal/logger"
)

type DBClient struct {
	db *sql.DB
}

func New(ctx context.Context, dns string) (*DBClient, error) {
	database, err := sql.Open("pgx", dns)
	if err != nil {
		logger.Log.Error("[DBClient] Could no create DB connection. Error: %s", err)
		return nil, err
	}

	databaseClient := &DBClient{
		db: database,
	}

	return databaseClient, nil
}

func (d *DBClient) Execute(ctx context.Context, query string, args ...any) error {
	var err error
	_, err = d.db.ExecContext(ctx, query, args...)

	if err != nil {
		logger.Log.Error("[DBClient] Could execute query %s. Error: %s", query, err)
		return err
	}
	return nil
}

func (d *DBClient) RowQuery(ctx context.Context, query string, args ...any) (*sql.Row, error) {
	resRow := d.db.QueryRowContext(ctx, query, args...)
	var err = resRow.Err()
	if err != nil {
		logger.Log.Error("[DBClient] Could execute query %s. Error: %s", query, err)
		return nil, err
	}
	return resRow, nil
}

func (d *DBClient) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	resRows, err := d.db.QueryContext(ctx, query, args...)
	if err != nil {
		logger.Log.Error("[DBClient] Could execute query %s. Error: %s", query, err)
		return nil, err
	}
	return resRows, nil
}

func (d *DBClient) Close() {
	if err := d.db.Close(); err != nil {
		logger.Log.Error("[DBClient] Could not close db connection. Error:  %s", err)
	}
}

func (d *DBClient) Ping(ctx context.Context) (bool, error) {
	if err := d.db.PingContext(ctx); err != nil {
		logger.Log.Error("[DBClient] Ping. Error: %s", err)
		return false, err
	}
	return true, nil
}

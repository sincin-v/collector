package db

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sincin-v/collector/internal/logger"
)

type DBClient struct {
	db *sql.DB
}

func New(dns string) (*DBClient, error) {
	database, err := sql.Open("pgx", dns)
	if err != nil {
		logger.Log.Error("Could no create DB connection. Error: %s", err)
		return nil, err
	}
	return &DBClient{
		db: database,
	}, nil
}

func (d *DBClient) Close() {
	if err := d.db.Close(); err != nil {
		logger.Log.Error("Could not close db connection. Error:  %s", err)
	}
}

func (d *DBClient) Ping(ctx context.Context) (bool, error) {
	ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	if err := d.db.PingContext(ctxTimeout); err != nil {
		logger.Log.Error("DB.Ping.Error: %s", err)
		return false, err
	}
	return true, nil
}

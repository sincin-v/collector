package storage

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sincin-v/collector/internal/logger"
	"github.com/sincin-v/collector/internal/models"
	"github.com/sincin-v/collector/internal/server/clients/db"
	"github.com/sincin-v/collector/internal/server/config"
)

type DBStorage struct {
	dbClient db.DBClient
}

func NewDBStorage(ctx context.Context, dbClient db.DBClient) (*DBStorage, error) {
	ds := DBStorage{
		dbClient: dbClient,
	}
	if err := ds.initialDataBaseTable(ctx); err != nil {
		logger.Log.Error("[DBClient] Could not init DBStorage. Error: %s", err)
		return nil, err
	}
	return &ds, nil
}

func (ds *DBStorage) initialDataBaseTable(ctx context.Context) error {
	query := "CREATE TABLE IF NOT EXISTS metrics (id SERIAL PRIMARY KEY, name VARCHAR UNIQUE, m_type VARCHAR, gauge_value DOUBLE PRECISION, counter_value BIGINT, CONSTRAINT UC_metric_name UNIQUE (name, m_type) );"
	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := ds.dbClient.Execute(ctxTimeout, query); err != nil {
		logger.Log.Error("[DBStorage] Could not create table 'metrics'. Error: %s", err)
		return err
	}
	return nil
}

func (ds *DBStorage) UpdateMetricsByBatch(ctx context.Context, metrics []models.Metrics) error {
	tx, err := ds.dbClient.BeginTx(ctx)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("[DBStorage] Cannot start Tx. Error: %s", err))
		return err
	}
	for _, m := range metrics {
		var errExec error
		switch m.MType {
		case config.CounterMetricType:
			errExec = ds.UpdateCounterMetric(ctx, m.ID, *m.Delta)
		case config.GaugeMetricType:
			errExec = ds.UpdateGaugeMetric(ctx, m.ID, *m.Value)
		default:
			logger.Log.Warnf("[DBStorage] Invalid type %s", m.MType)
			continue
		}
		if errExec != nil {
			err = errors.Join(err, fmt.Errorf("[DBStorage] Execute query error: %s", errExec))
			errRollback := tx.Rollback()
			if errRollback != nil {
				err = errors.Join(err, fmt.Errorf("[DBStorage] Rollback query error: %s", errRollback))
			}
			return err
		}
	}
	errCommit := tx.Commit()
	if errCommit != nil {
		err = errors.Join(err, fmt.Errorf("[DBStorage] Commit error %s", errCommit))
		return err
	}
	return nil
}

func (ds *DBStorage) UpdateGaugeMetric(ctx context.Context, name string, value float64) error {
	query := "INSERT INTO metrics (name, m_type, gauge_value, counter_value) VALUES ($1, 'gauge', $2, NULL) ON CONFLICT (name) DO UPDATE SET gauge_value = excluded.gauge_value"
	if err := ds.dbClient.Execute(ctx, query, name, value); err != nil {
		err = errors.Join(err, fmt.Errorf("[DBStorage] Could not update metric %s. Error: %s", name, err))
		return err
	}
	return nil
}

func (ds *DBStorage) UpdateCounterMetric(ctx context.Context, name string, value int64) error {
	query := "INSERT INTO metrics (name, m_type, gauge_value, counter_value) VALUES ($1, 'counter', NULL, $2) ON CONFLICT (name) DO UPDATE SET counter_value = metrics.counter_value+EXCLUDED.counter_value"
	if err := ds.dbClient.Execute(ctx, query, name, value); err != nil {
		err = errors.Join(err, fmt.Errorf("[DBStorage] Could not update metric %s. Error: %s", name, err))
		return err
	}
	return nil
}

func (ds *DBStorage) GetMetric(ctx context.Context, metricType string, metricName string) (string, error) {
	query := "SELECT * FROM metrics WHERE name=$1 AND m_type=$2"
	var err error
	resRow, err := ds.dbClient.RowQuery(ctx, query, metricName, metricType)
	if err != nil {
		err = errors.Join(err, fmt.Errorf("[DBStorage] metric %s does not exists Error: %s", metricName, err))
		return "", err
	}
	var m models.DBMetricModel
	errScan := resRow.Scan(&m.ID, &m.Name, &m.MType, &m.GaugeValue, &m.CounterValue)
	if errScan != nil {
		err = errors.Join(err, fmt.Errorf("[DBStorage] metric %s does not exists Error: %s", metricName, errScan))
		return "", err
	}
	switch metricType {
	case config.GaugeMetricType:
		return strconv.FormatFloat(*m.GaugeValue, 'f', -1, 64), nil
	case config.CounterMetricType:
		return fmt.Sprintf("%d", *m.CounterValue), nil
	default:
		return "", err
	}
}

func (ds *DBStorage) GetAllCountersMetrics(ctx context.Context) map[string]int64 {
	query := "SELECT * FROM metrics WHERE m_type = 'counter'"
	var result = make(map[string]int64)
	resRows, err := ds.dbClient.Query(ctx, query)
	if err != nil {
		logger.Log.Error("[DBStorage] Could get counter metrics Error: %s", err)
		return nil
	}
	for resRows.Next() {
		var m models.DBMetricModel
		errScan := resRows.Scan(&m.ID, &m.Name, &m.MType, &m.GaugeValue, &m.CounterValue)
		if errScan != nil {
			logger.Log.Errorf("[DBStorage] Cannot get scan counter metrics object. Error: %s", errScan)
			return nil
		}
		result[*m.Name] = *m.CounterValue
	}
	defer func() {
		if errClose := resRows.Close(); errClose != nil {
			err = errors.Join(err, fmt.Errorf("[DBStorage] close rows error: %w", errClose))
		}
	}()
	return result
}

func (ds *DBStorage) GetAllGaugeMetrics(ctx context.Context) map[string]float64 {
	query := "SELECT * FROM metrics WHERE m_type = 'gauge'"
	var result = make(map[string]float64)
	resRows, err := ds.dbClient.Query(ctx, query)
	if err != nil {
		logger.Log.Error("[DBStorage] Could get gauge metrics Error: %s", err)
		return nil
	}
	for resRows.Next() {
		var m models.DBMetricModel
		errScan := resRows.Scan(&m.ID, &m.Name, &m.MType, &m.GaugeValue, &m.CounterValue)
		if errScan != nil {
			logger.Log.Errorf("[DBStorage] Cannot get scan gauge metrics object. Error: %s", errScan)
			return nil
		}
		result[*m.Name] = *m.GaugeValue
	}
	return result
}

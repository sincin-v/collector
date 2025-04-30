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
)

type DBStorage struct {
	baseCtx  context.Context
	dbClient db.DBClient
}

func NewDBStorage(ctx context.Context, dbClient db.DBClient) (*DBStorage, error) {
	ds := DBStorage{
		baseCtx:  ctx,
		dbClient: dbClient,
	}
	if err := ds.initialDataBaseTable(ctx); err != nil {
		logger.Log.Error("[DBClient] Could not init DBStorage. Error: %s", err)
		return nil, err
	}
	return &ds, nil
}

func (ds *DBStorage) initialDataBaseTable(ctx context.Context) error {
	query := "CREATE TABLE IF NOT EXISTS metrics (id SERIAL PRIMARY KEY, name VARCHAR UNIQUE, m_type VARCHAR, gauge_value DOUBLE PRECISION, counter_value INTEGER, CONSTRAINT UC_metric_name UNIQUE (name) );"
	ctxTimeout, cancel := context.WithTimeout(ds.baseCtx, 5*time.Second)
	defer cancel()
	if err := ds.dbClient.Execute(ctxTimeout, query); err != nil {
		logger.Log.Error("[DBStorage] Could not create table 'metrics'. Error: %s", err)
		return err
	}
	return nil
}

func (ds *DBStorage) UpdateGaugeMetric(name string, value float64) {
	query := "INSERT INTO metrics (name, m_type, gauge_value, counter_value) VALUES ($1, 'gauge', $2, NULL) ON CONFLICT (name) DO UPDATE SET gauge_value = excluded.gauge_value"
	ctxTimeout, cancel := context.WithTimeout(ds.baseCtx, 5*time.Second)
	defer cancel()
	if err := ds.dbClient.Execute(ctxTimeout, query, name, value); err != nil {
		logger.Log.Errorf("[DBStorage] Could not update metric %s. Error: %s", name, err)
	}
}

func (ds *DBStorage) UpdateCounterMetric(name string, value int64) {
	query := "INSERT INTO metrics (name, m_type, gauge_value, counter_value) VALUES ($1, 'counter', NULL, $2) ON CONFLICT (name) DO UPDATE SET counter_value = metrics.counter_value+EXCLUDED.counter_value"
	ctxTimeout, cancel := context.WithTimeout(ds.baseCtx, 5*time.Second)
	defer cancel()
	if err := ds.dbClient.Execute(ctxTimeout, query, name, value); err != nil {
		logger.Log.Errorf("[DBStorage] Could not update metric %s. Error: %s", name, err)
	}

}

func (ds *DBStorage) GetMetric(metricType string, metricName string) (string, error) {
	query := "SELECT * FROM metrics WHERE name=$1 AND m_type=$2"
	ctxTimeout, cancel := context.WithTimeout(ds.baseCtx, 5*time.Second)
	defer cancel()
	resRow, err := ds.dbClient.RowQuery(ctxTimeout, query, metricName, metricType)
	if err != nil {
		logger.Log.Error("[DBStorage] Could get  metric %s. Error: %s", metricName, err)
		return "", err
	}
	var m models.DBMetricModel
	errScan := resRow.Scan(&m.ID, &m.Name, &m.MType, &m.GaugeValue, &m.CounterValue)
	if errScan != nil {
		logger.Log.Errorf("[DBStorage] Cannot get scan metric %s object. Error: %s", metricName, errScan)
		return "", errScan
	}
	switch metricType {
	case "gauge":
		return strconv.FormatFloat(*m.GaugeValue, 'f', -1, 64), nil
	case "counter":
		return fmt.Sprintf("%d", *m.CounterValue), nil
	default:
		return "", errors.New("[DBStorage] Invalid metric type")
	}
}

func (ds *DBStorage) GetAllCountersMetrics() map[string]int64 {
	query := "SELECT * FROM metrics WHERE m_type = 'counter'"
	ctxTimeout, cancel := context.WithTimeout(ds.baseCtx, 5*time.Second)
	defer cancel()
	var result = make(map[string]int64)
	resRows, err := ds.dbClient.Query(ctxTimeout, query)
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
	return result
}

func (ds *DBStorage) GetAllGaugeMetrics() map[string]float64 {
	query := "SELECT * FROM metrics WHERE m_type = 'gauge'"
	ctxTimeout, cancel := context.WithTimeout(ds.baseCtx, 5*time.Second)
	defer cancel()
	var result = make(map[string]float64)
	resRows, err := ds.dbClient.Query(ctxTimeout, query)
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

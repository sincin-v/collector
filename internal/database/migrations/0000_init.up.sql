CREATE TABLE IF NOT EXISTS metrics (id SERIAL PRIMARY KEY, name VARCHAR UNIQUE, m_type VARCHAR, gauge_value DOUBLE PRECISION, counter_value BIGINT, CONSTRAINT UC_metric_name UNIQUE (name, m_type) );

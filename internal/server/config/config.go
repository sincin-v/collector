package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host            string `env:"ADDRESS"`
	StoreInterval   int64  `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	LogLevel        string `env:"LOG_LEVEL" envDefault:"INFO"`
	DBDns           string `env:"DATABASE_DSN"`

	RetryIntervals []time.Duration `env:"RETRY_INTERVALS" envSeparator:"," envDefault:"1s,3s,5s"`
}

func GetServerConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Host, "a", "localhost:8080", "Listen host and port")
	flag.Int64Var(&cfg.StoreInterval, "i", 300, "interval store metrics value")
	flag.StringVar(&cfg.FileStoragePath, "f", "/tmp/metric_storage", "Path to storage file")
	flag.BoolVar(&cfg.Restore, "r", true, "Flag of restore collected metrics data")
	flag.StringVar(&cfg.DBDns, "d", "", "DNS fo connect to DB")
	flag.Parse()
	var err = env.Parse(cfg)
	if err != nil {
		log.Printf("Cannot parse env")
		return nil, err
	}
	return cfg, nil
}

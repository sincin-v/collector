package config

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host     string `env:"ADDRESS"`
	LogLevel string `env:"LOG_LEVEL" envDefault:"INFO"`
}

func GetServerConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.Host, "a", "localhost:8080", "Listen host and port")
	flag.Parse()
	var err = env.Parse(cfg)
	if err != nil {
		log.Printf("Cannot parse env")
		return nil, err
	}
	return cfg, nil
}

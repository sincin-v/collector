package config

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerHost     string        `env:"ADDRESS"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
}

type AgentConfig struct {
	ServerHost     string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func GetAgentConfig() (*Config, error) {
	cfg := &Config{}

	flag.StringVar(&cfg.ServerHost, "a", "localhost:8080", "Metric server host and port")
	var argReportInterval = flag.Int("r", 10, "Report interval")
	var argPollInterval = flag.Int("p", 2, "Poll interval")

	flag.Parse()

	err := env.Parse(&cfg)
	if err != nil {
		log.Printf("Cannot parse env")
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = time.Duration(time.Duration(*argReportInterval) * time.Second)
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = time.Duration(time.Duration(*argPollInterval) * time.Second)

	}

	return cfg, nil
}

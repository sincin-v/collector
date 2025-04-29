package logger

import "go.uber.org/zap"

var Log zap.SugaredLogger

func Initialize(level string) error {
	logLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = logLevel

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = *logger.Sugar()
	return nil
}

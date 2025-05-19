package migrator

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sincin-v/collector/internal/logger"
)


func ApplyMigrations(databaseDSN string, migrationPath string) error {
	m, err := migrate.New(
		"file://" +  migrationPath,
		databaseDSN,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			logger.Log.Debug("no migrations to apply")
			return nil
		}
		return err
	}
	logger.Log.Debug("migrations applied successfully")
	return nil
}

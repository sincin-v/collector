package migrator

import (
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
)


func ApplyMigrations(databaseDSN string, migrationPath string) error {
	m, err := migrate.New(
		"file://" + migrationPath,
		databaseDSN,
	)

	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		return err
	}
	return nil
}

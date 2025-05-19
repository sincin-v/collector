package main

import (
	"flag"
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/sincin-v/collector/internal/database/migrator"
)

func main() {

	var migrationsPath, databaseDSN string

	flag.StringVar(
		&databaseDSN,
		"d",
		"",
		"database-dsn",
	)

	flag.StringVar(
		&migrationsPath,
		"p",
		"",
		"path to migrations",
	)
	flag.Parse()
	if migrationsPath == "" {
		panic("migrations path is required")
	}

	migrateErr := migrator.ApplyMigrations(databaseDSN, migrationsPath)

	if migrateErr != nil {
		panic(migrateErr)
	}

	fmt.Println("migrations applied successfully")
}

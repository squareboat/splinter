package runner

import (
	"context"
	"log"

	"github.com/the-e3n/splinter/database/postgres"
	"github.com/the-e3n/splinter/logger"
	"github.com/the-e3n/splinter/parser"
)

func Postgres(connURL, migrationType string) {
	if migrationType == "" {
		logger.Log.Error("invalid migration type")
		return
	}

	ctx := context.TODO()
	driver, err := postgres.NewPostgresDB(connURL, migrationType)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = driver.Initialize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// place locks
	err = driver.Lock()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err = driver.Unlock()
		if err != nil {
			logger.Log.Error(err)
		}
	}()
	// get migration files
	migrationFiles, err := parser.GetMigrationFileNames()
	if err != nil {
		log.Fatal(err)
	}

	migrationsToExec, err := driver.CrossCheckMigrations(ctx, migrationFiles)
	if err != nil {
		logger.Log.Error(err)
		return
	}

	if len(migrationsToExec) > 0 {
		err = driver.Migrate(ctx, migrationsToExec)
		if err != nil {
			logger.Log.WithError(err)
		}
	}

}

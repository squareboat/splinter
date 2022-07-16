package runner

import (
	"context"
	"log"

	"github.com/the-e3n/splinter/database/postgres"
	"github.com/the-e3n/splinter/parser"
)

func Postgres(connURL, migrationType string) {
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

	// get migration files
	migrationFiles, _ := parser.GetMigrationFileNames()
	newMigrations, err := driver.CrossCheckMigrations(ctx, migrationFiles)
	if err != nil {
		log.Fatal(err)
	}
	if len(newMigrations) > 0 {
		err = driver.Migrate(ctx, newMigrations)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = driver.Unlock()
	if err != nil {
		log.Fatal(err)
	}
}

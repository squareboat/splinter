package runner

import (
	"context"
	"log"

	"github.com/the-e3n/splinter/database/postgres"
)

func Postgres(connURL, migrationType string) {

	driver, err := postgres.NewPostgresDB(connURL, migrationType)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = driver.Initialize(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// get migration files

}

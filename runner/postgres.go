package runner

import (
	"context"
	"log"

	"github.com/squareboat/splinter/constants"
	"github.com/squareboat/splinter/database/postgres"
	"github.com/squareboat/splinter/logger"
)

func Postgres(connURL, migrationType string) {
	// if migrationType == "" {
	// 	logger.Log.Error("invalid migration type")
	// 	return
	// }

	// ctx := context.TODO()
	// driver, err := postgres.NewPostgresDB(connURL, migrationType)
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// err = driver.Initialize(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // place locks
	// err = driver.Lock()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer func() {
	// 	err = driver.Unlock()
	// 	if err != nil {
	// 		logger.Log.Error(err)
	// 	}
	// }()
	// // get migration files
	// migrationFiles, err := parser.GetMigrationFileNames()
	// if err != nil {
	// 	logger.Log.Error(err)
	// 	return
	// }

	// migrationsToExec, err := driver.CrossCheckMigrations(ctx, migrationFiles)
	// if err != nil {
	// 	logger.Log.Error(err)
	// 	return
	// }

	// if len(migrationsToExec) == 0 {
	// 	if migrationType == constants.MIGRATION_UP {
	// 		logger.Log.Warn("No new migrations found.")
	// 		return
	// 	}

	// 	logger.Log.Warn("No migrations found")
	// }
	// if len(migrationsToExec) > 0 {
	// 	err = driver.Migrate(ctx, migrationsToExec)
	// 	if err != nil {
	// 		logger.Log.WithError(err)
	// 	}
	// }

}
func UnlockDB(connURL string) {
	ctx := context.TODO()
	driver, err := postgres.NewPostgresDB(connURL, constants.UNLOCK_DB)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = driver.Initialize(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Unlock DB
	logger.Log.Info("Unlocking Database ...")
	err = driver.Unlock()
	if err != nil {
		logger.Log.Error(err)
	}

}

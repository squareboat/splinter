package migrate

import (
	"context"
	"errors"
	"log"
	"sort"
	"time"

	"github.com/squareboat/splinter/config"
	"github.com/squareboat/splinter/constants"
	"github.com/squareboat/splinter/database"
	"github.com/squareboat/splinter/database/postgres"
	"github.com/squareboat/splinter/logger"
	"github.com/squareboat/splinter/parser"
)

type Migrate struct {
	driver database.Driver
	dbUri  string
}

func NewMigrate(migrationType string) (*Migrate, error) {
	driver, err := postgres.NewPostgresDriver(config.GetDbUri(), migrationType)
	if err != nil {
		logger.Log.Error(err)
		log.Fatal(err)
	}
	migrate := Migrate{
		driver: driver,
		dbUri:  config.GetDbUri(),
	}
	return &migrate, nil
}

func (m *Migrate) Up() error {
	err := m.driver.Initialize(context.Background())
	if err != nil {
		logger.Log.WithError(err)
		return err
	}
	err = m.lock()
	if err != nil {
		logger.Log.WithError(err)
		return err
	}

	defer func() {
		m.unlock()
		m.close()
	}()

	migrations, err := m.driver.GetSchemaMigrations()
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	migrationFiles, err := parser.GetMigrationFileNames()
	if err != nil {
		logger.Log.Error(err)
		return err
	}
	logger.Log.Info("Migration from files ", migrationFiles)
	logger.Log.Info("Schema Migations ", migrations)
	err = m.driver.CrossCheckMigrations(context.Background(), migrationFiles, migrations)
	if err != nil {
		logger.Log.Warn("error cross checking migrations")
		return err
	}
	logger.Log.Info("cross check success")
	if len(migrationFiles) == len(migrations) {
		logger.Log.Warn("no new migrations found")
		return nil
	}

	newMigrations := getNewMigrations(migrationFiles, migrations)
	logger.Log.Info("Migrations  ", newMigrations)

	//migrationFiles[len(migrations):]
	err = m.driver.Migrate(context.Background(), newMigrations)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	lastBatchNumber := 1
	if len(migrations) > 0 {
		lastBatchNumber = migrations[len(migrations)-1].BatchNumber + 1
	}

	newSchemaMigrations := []database.SchemaMigration{}
	for i := range newMigrations {
		newSchemaMigrations = append(newSchemaMigrations, database.SchemaMigration{
			MigrationName: newMigrations[i],
			BatchNumber:   lastBatchNumber,
			CreatedAt:     time.Now().Unix(),
		})
	}

	// insert into schema_migations
	err = m.driver.UpdateSchemaMigrations(newSchemaMigrations, constants.MIGRATION_UP)
	if err != nil {
		logger.Log.WithError(err).Error("error updating schema_migrations")
		return err
	}
	return nil
}

func (m *Migrate) Down() error {
	err := m.driver.Initialize(context.Background())
	if err != nil {
		logger.Log.WithError(err)
		return err
	}
	err = m.lock()
	if err != nil {
		logger.Log.WithError(err)
		return err
	}

	defer func() {
		m.unlock()
		m.close()
	}()

	migrations, err := m.driver.GetSchemaMigrations()
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	migrationFiles, err := parser.GetMigrationFileNames()
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	err = m.driver.CrossCheckMigrations(context.Background(), migrationFiles, migrations)
	if err != nil {
		logger.Log.Warn("error cross checking migrations")
		return err
	}
	if len(migrations) == 0 {
		return errors.New("no migrations found")
	}

	rollbackMigrations := getRollbacks(1, migrations)

	err = m.driver.Migrate(context.Background(), rollbackMigrations)
	if err != nil {
		logger.Log.Error(err)
		return err
	}
	schemaMigrations := []database.SchemaMigration{}
	for i := range rollbackMigrations {
		schemaMigrations = append(schemaMigrations, database.SchemaMigration{MigrationName: rollbackMigrations[i]})
	}
	err = m.driver.UpdateSchemaMigrations(schemaMigrations, constants.MIGRATION_DOWN)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	return nil
}

func (m *Migrate) lock() error {
	return m.driver.Lock()
}

func (m *Migrate) unlock() error {
	return m.driver.Unlock()
}

func (m *Migrate) close() error {
	return m.driver.Close()
}

func getNewMigrations(migrationFiles []string, schemaMigrations []database.SchemaMigration) []string {
	newMigrations := []string{}
	migrationsFilesMap := map[string]bool{}
	for i := range schemaMigrations {
		migrationsFilesMap[schemaMigrations[i].MigrationName] = true
	}

	// order of files is guaranteed to be in sorted order
	for i := range migrationFiles {
		if _, ok := migrationsFilesMap[migrationFiles[i]]; !ok {
			newMigrations = append(newMigrations, migrationFiles[i])
		}
	}
	sort.Slice(newMigrations, func(i, j int) bool {
		return newMigrations[i] < newMigrations[j]
	})
	return newMigrations
}

func getRollbacks(count int, schemaMigrations []database.SchemaMigration) []string {
	return []string{schemaMigrations[len(schemaMigrations)-1].MigrationName}
}

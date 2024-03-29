package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/squareboat/splinter/constants"
	"github.com/squareboat/splinter/database"
	"github.com/squareboat/splinter/logger"
	"github.com/squareboat/splinter/parser"
)

type Postgres struct {
	// conn *sql.Conn
	db                *sql.DB
	dbName            string
	migrationType     string
	latestBatchNumber int
}

func (p *Postgres) Initialize(ctx context.Context) error {
	// check if schema migrations table is present or not

	query := tableExists(constants.SCHEMA_MIGRATIONS, constants.DEFAULT_SCHEMA_NAME)

	res, err := p.db.Query(query)
	if err != nil {
		logger.Log.WithError(err)
		return err
	}

	for res.Next() {
		var row interface{}
		res.Scan(&row)

		tableExists, ok := row.(bool)

		if !ok {
			return errors.New("error reading schema_migrations table")
		}

		if !tableExists {

			if _, err := p.db.Exec(createSchemaMigrations()); err != nil {
				logger.Log.WithError(err).Error("error creating migration table")
				return err
			}

		}
	}

	// check if migrations lock table is present
	migrationsLockQuery := tableExists(constants.MIGRTATION_LOCKS, constants.DEFAULT_SCHEMA_NAME)
	migrationLockResponse, err := p.db.Query(migrationsLockQuery)
	if err != nil {
		return err
	}

	for migrationLockResponse.Next() {
		var row interface{}
		migrationLockResponse.Scan(&row)
		fmt.Println("migration response ", row)
		if tableExists, ok := row.(bool); ok && !tableExists {
			// if !tableExists {

			_, err = p.db.Exec(createMigrationLocksTable())
			if err != nil {
				logger.Log.WithError(err).Error("error creating migrations lock table")
				return err
			}

			// insert lock row
			_, err = p.db.Exec(insertMigrationLock())
			if err != nil {
				logger.Log.WithError(err)
				return err
			}

			// }
		}

	}

	return nil
}

// func (p *Postgres) CrossCheckMigrations(ctx context.Context, migrationFiles []string) ([]string, error) {
// 	// read from schema_migrations
// 	// case 1 migration file does not exist in the table, then execute those migrations
// 	// case 2 migration exists in database but does not exist in file system, then throw an error and mark the migration as dirty.

// 	query := getMigrations()
// 	isNewMigration := map[string]bool{}
// 	maxBatchNumner := 0
// 	migrations := []schemaMigration{}
// 	latestMigrationFile := []string{}
// 	firstIteration := true
// 	newMigrations := []string{}

// 	for i := range migrationFiles {
// 		isNewMigration[migrationFiles[i]] = true
// 	}
// 	sqlRows, err := p.db.Query(query)
// 	if err != nil {
// 		logger.Log.WithError(err)
// 		return nil, err
// 	}

// 	for sqlRows.Next() {
// 		var (
// 			id            int64
// 			migrationName string
// 			batchNumber   int
// 			createdAt     int64
// 		)

// 		if err = sqlRows.Scan(&id, &migrationName, &batchNumber, &createdAt); err != nil {
// 			logger.Log.WithError(err)
// 			return nil, err
// 		}

// 		if batchNumber > maxBatchNumner {
// 			maxBatchNumner = batchNumber
// 		}
// 		mig := schemaMigration{
// 			migrationName: migrationName,
// 			id:            id,
// 			createdAt:     createdAt,
// 			batchNumber:   batchNumber,
// 		}
// 		migrations = append(migrations, mig)

// 		if firstIteration {
// 			firstIteration = false
// 			latestMigrationFile = append(latestMigrationFile, mig.migrationName)
// 		}

// 	}
// 	p.latestBatchNumber = maxBatchNumner

// 	for i := range migrations {
// 		migrationFromDB := migrations[i]
// 		if _, ok := isNewMigration[migrationFromDB.migrationName]; !ok {
// 			return nil, errors.New("some migration files  missing in your migration path")
// 		}
// 		isNewMigration[migrationFromDB.migrationName] = false
// 	}

// 	for fileName, isNew := range isNewMigration {
// 		if isNew {
// 			newMigrations = append(newMigrations, fileName)
// 		}
// 	}

// 	if p.migrationType == constants.MIGRATION_DOWN {
// 		return latestMigrationFile, nil
// 	}

// 	sort.Slice(newMigrations, func(i, j int) bool {
// 		return newMigrations[j] > newMigrations[i]

// 	})
// 	return newMigrations, nil
// }

// TODO implment CrossCheckMigrations
func (p *Postgres) CrossCheckMigrations(ctx context.Context, migrations []string, schemaMigrations []database.SchemaMigration) error {
	if len(migrations) < len(schemaMigrations) {
		return errors.New("missing migration files")
	}

	migrationsFilesMap := map[string]bool{}
	for i := range migrations {
		migrationsFilesMap[migrations[i]] = true
	}

	// order of files is guaranteed to be in sorted order
	for i := range schemaMigrations {
		if _, ok := migrationsFilesMap[schemaMigrations[i].MigrationName]; !ok {
			return fmt.Errorf("missing migration file %v", schemaMigrations[i].MigrationName)
		}
	}

	return nil
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

func (p *Postgres) GetSchemaMigrations() ([]database.SchemaMigration, error) {
	query := getMigrations()
	result := []database.SchemaMigration{}
	sqlRows, err := p.db.Query(query)
	if err != nil {
		logger.Log.WithError(err)
		return nil, err
	}

	for sqlRows.Next() {
		var (
			id            int64
			migrationName string
			batchNumber   int
			createdAt     int64
		)

		if err = sqlRows.Scan(&id, &migrationName, &batchNumber, &createdAt); err != nil {
			logger.Log.WithError(err)
			return nil, err
		}

		mig := database.SchemaMigration{
			MigrationName: migrationName,
			ID:            id,
			CreatedAt:     createdAt,
			BatchNumber:   batchNumber,
		}
		result = append(result, mig)

	}

	return result, nil

}

func (p *Postgres) UpdateSchemaMigrations(migrations []database.SchemaMigration, migrationType string) error {
	debug := "error updaing schema migrations"

	if migrationType == constants.MIGRATION_UP {
		query := updateSchemaMigrations(migrations)
		err := p.execQuery(query)
		if err != nil {
			logger.Log.WithError(err).Error(debug)
			return err
		}
		return nil
	}

	if migrationType == constants.MIGRATION_DOWN {
		query := deleteFromSchemaMigrations(migrations)
		logger.Log.Info("query ", query)
		err := p.execQuery(query)
		if err != nil {
			logger.Log.WithError(err).Error(debug)
			return err
		}
		return nil
	}

	return errors.New("invalid migrtion type")
}

func (p *Postgres) execQuery(q string) error {
	_, err := p.db.Exec(q)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	return nil

}

// runs given set of SQL
func (p *Postgres) Migrate(ctx context.Context, migrationFiles []string) error {
	transaction, err := p.db.BeginTx(ctx, nil)

	if err != nil {
		logger.Log.WithError(err).Error("error begin transaction")
		return err
	}

	for _, filename := range migrationFiles {

		// get queries in the file
		queries, err := parser.ParseFile(filename, p.migrationType)
		if err != nil {
			logger.Log.Errorf("error parsing file %s", filename)
			logger.Log.Error(err)
			transaction.Rollback()
			return err
		}

		for i := range queries {
			q := queries[i]

			_, err = transaction.Exec(q)
			if err != nil {

				logger.Log.Warn("rolling back transaction")
				logger.Log.Error(err)
				logger.Log.Error("Migration file ", filename, " \nQuery : ", q)
				rollbackErr := transaction.Rollback()

				if rollbackErr != nil {
					logger.Log.WithError(rollbackErr).Error("error rolling back")
				}

				return err
			}
		}

	}

	logger.Log.Info("Commiting transaction")
	err = transaction.Commit()
	if err != nil {
		logger.Log.WithError(err)
		return err
	}
	return nil
}

func (p *Postgres) updateSchemaMigrations(migrationFiles []string) string {

	if p.migrationType == constants.MIGRATION_UP {
		query := insertSchemaMigrations(migrationFiles, p.latestBatchNumber+1)
		logger.Log.Info("update ", query)
		return query
	}

	if p.migrationType == constants.MIGRATION_DOWN {
		return deleteLatestSchemaMigrations(migrationFiles[0])
	}

	return ""
}

func (p *Postgres) Validate(migrations []string) error {
	return nil
}

func (p *Postgres) Lock() error {
	// check if migration is unlocked
	query := getLock(true)
	sqlRows, err := p.db.Query(query)
	if err != nil {
		logger.Log.WithError(err)
		return err
	}
	// TODO: handle if now rows found in migration locks table

	for sqlRows.Next() {
		var (
			id       int
			isLocked bool
		)

		err = sqlRows.Scan(&id, &isLocked)
		if err != nil {
			logger.Log.WithError(err)
			return err
		}

		if isLocked {
			logger.Log.Warn("migration table is already locked")
			return errors.New("can't take lock to run migrations: migration table is already locked")
		}
	}
	// set is lock to true

	query = updateMigrationLock(true)
	sqlRes, err := p.db.Exec(query)
	if err != nil {
		logger.Log.Error(err)
		return err
	}

	rowsAffected, err := sqlRes.RowsAffected()
	if err != nil {
		logger.Log.Warn(err)
	}
	fmt.Println("rows affected", rowsAffected)
	return nil
}

func (p *Postgres) Unlock() error {
	// if lock is not present throw error

	query := updateMigrationLock(false)
	sqlRes, err := p.db.Exec(query)
	if err != nil {
		logger.Log.WithError(err).Error()
		return err
	}
	rowsAffected, err := sqlRes.RowsAffected()

	if err != nil {
		logger.Log.WithError(err).Error()
		return err
	}

	if rowsAffected == 0 {
		return errors.New("unable to remove lock. no locks found")
	}

	logger.Log.Info("migration lock removed successfully")
	return nil
}

func NewPostgresDB(connectionURL, migrationType string) (database.Driver, error) {
	// parse connection URL
	driver := Postgres{}
	db, err := sql.Open("postgres", connectionURL)

	if err != nil {
		logger.Log.WithError(err)
		panic(err)
	}

	if migrationType == "" {
		return nil, errors.New("invalid migration type")
	}
	var dbName string

	// get dbname from connection url
	connURLParts := strings.Split(connectionURL, "/")
	// db_name is the last element of the array formed
	if len(connURLParts) > 0 {

		// to get rid of query params
		dbNameParts := strings.Split(connURLParts[len(connURLParts)-1], "?")
		if len(dbNameParts) > 0 {
			name := dbNameParts[0]
			dbName = name
		}
	}

	if dbName == "" {
		return nil, errors.New("invalid database name")
	}

	driver.db = db
	driver.dbName = dbName
	driver.migrationType = migrationType

	return &driver, nil
}

func NewPostgresDriver(connectionURL string, migrationType string) (database.Driver, error) {
	db, err := sql.Open("postgres", connectionURL)

	if err != nil {
		logger.Log.WithError(err)
		return nil, err
	}
	return &Postgres{db: db, migrationType: migrationType}, nil
}

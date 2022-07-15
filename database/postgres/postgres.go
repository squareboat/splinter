package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq"
	"github.com/the-e3n/splinter/constants"
	"github.com/the-e3n/splinter/database"
	"github.com/the-e3n/splinter/logger"
)

type Postgres struct {
	// conn *sql.Conn
	db            *sql.DB
	dbName        string
	migrationType string
}

func (p *Postgres) Initialize(ctx context.Context) error {
	// check if schema migrations table is present or not

	query := tableExists(constants.SCHEMA_MIGRATIONS, constants.DEFAULT_SCHEMA_NAME)

	res, err := p.db.Query(query)
	if err != nil {
		logger.Log.WithError(err)
	}

	for res.Next() {
		var row interface{}
		res.Scan(&row)

		tableExists, ok := row.(bool)
		if ok {
			// create table schema migrations
			if !tableExists {

				if _, err := p.db.Exec(createSchemaMigrations()); err != nil {
					logger.Log.WithError(err).Error("error creating migration table")
					return err
				}

			}

		} else {
			return errors.New("error reading schema_migrations table")
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
		if tableExists, ok := row.(bool); ok {
			if !tableExists {

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

			}
		}

	}

	return nil
}

func (p *Postgres) CrossCheckMigrations(ctx context.Context, migrationFiles []string) (map[string]string, error) {
	// read from schema_migrations
	// case 1 migration file does not exist in the table, then execute those migrations
	// case 2 migration exists in database but does not exist in file system, then throw an error and mark the migration as dirty.
	fmt.Println("cross checking migration files from DBc")
	query := getMigrations()
	isNewMigration := map[string]bool{}
	for i := range migrationFiles {
		isNewMigration[migrationFiles[i]] = true
	}
	sqlRows, err := p.db.Query(query)
	if err != nil {
		logger.Log.WithError(err)
		return nil, err
	}

	fmt.Println("sqlRows", sqlRows)

	migrations := []schemaMigration{}
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

		migrations = append(migrations, schemaMigration{
			migrationName: migrationName,
			id:            id,
			createdAt:     createdAt,
			batchNumber:   batchNumber,
		})

	}

	for i := range migrations {
		migrationFromDB := migrations[i]
		if _, ok := isNewMigration[migrationFromDB.migrationName]; !ok {
			log.Fatal("migration file from DB not found in yout path")
		}
		isNewMigration[migrationFromDB.migrationName] = false
	}

	return nil, nil
}

// runs given set of SQL
func (p *Postgres) Migrate(ctx context.Context, migrations map[string]string) error {
	transaction, err := p.db.BeginTx(ctx, nil)

	if err != nil {
		logger.Log.WithError(err)
		return err
	}

	for _, query := range migrations {
		_, err := transaction.Exec(query)

		if err != nil {
			logger.Log.WithError(err)
			rollbackErr := transaction.Rollback()

			if rollbackErr != nil {
				logger.Log.WithError(rollbackErr).Error("error rolling back")
			}

			return err
		}
	}

	err = transaction.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
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
			logger.Log.Fatal("Can't take lock to run migrations: Migration table is already locked")
		}
	}
	// set is lock to true

	query = updateMigrationLock(true)
	sqlRes, err := p.db.Exec(query)
	if err != nil {
		logger.Log.Fatal(err)
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

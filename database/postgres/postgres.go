package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/the-e3n/migrator/constants"
	"github.com/the-e3n/migrator/database"
	"github.com/the-e3n/migrator/logger"
)

type Postgres struct {
	// conn *sql.Conn
	db     *sql.DB
	dbName string
}

func (p *Postgres) Initialize(ctx context.Context) error {
	// check if schema migrations table is present or not

	query := tableExists(constants.SCHEMA_MIGRATIONS, constants.DEFAULT_SCHEMA_NAME)
	fmt.Println("query ", query)

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

		if tableExists, ok := row.(bool); ok {
			if !tableExists {
				_, err = p.db.Exec(createMigrationLocksTable())
				if err != nil {
					logger.Log.WithError(err).Error("error creating migrations lock table")
					return err
				}
			}
		}

	}

	return nil
}

func (p *Postgres) CrossCheckMigrations(ctx context.Context, migrationFiles []string) (map[string]string, error) {
	return nil, nil
}

// runs given set of SQL
func (p *Postgres) RunMigrations(ctx context.Context, migrations map[string]string) error {
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
		panic(err)
	}
	return nil
}

func (p *Postgres) Validate(migrations []string) error {
	return nil
}

func NewPostgresDB(connectionURL string) (database.Driver, error) {
	// parse connection URL
	driver := Postgres{}
	db, err := sql.Open("postgres", connectionURL)

	if err != nil {
		logger.Log.WithError(err)
		panic(err)
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

	return &driver, nil
}

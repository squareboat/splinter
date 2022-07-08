package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Postgres struct {
	// conn *sql.Conn
	db     *sql.DB
	dbName string
}

// runs given set of SQL
func (p *Postgres) RunMigration(ctx context.Context, migrations map[string]string) error {
	transaction, err := p.db.BeginTx(ctx, nil)

	if err != nil {
		fmt.Println("error 1", err)
		return err
	}

	for _, query := range migrations {
		result, err := transaction.Exec(query)

		if err != nil {
			fmt.Println("err", err)
			rollbackErr := transaction.Rollback()

			if rollbackErr != nil {
				fmt.Println("error rolling back", rollbackErr)
			}

			return err
		}
		fmt.Println("transaction result", result)
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

func NewPostgresDB(connectionURL string) (Driver, error) {
	// parse connection URL
	p := Postgres{}
	fmt.Println("url", connectionURL)
	db, err := sql.Open("postgres", connectionURL)

	if err != nil {
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
	p.db = db

	p.dbName = dbName

	return &p, nil
}

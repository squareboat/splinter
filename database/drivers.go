package database

import "context"

// a common driver interface that all db drivers will have to implement
type Driver interface {
	// Close connection to database
	Close() error
	// lock will update the migrations_lock table's entry
	// if lock is already present that means some other process is executing a migraion or
	// previos migrations errored out, either ways we cannot proceed with the migration
	Lock() error

	// unlock releases the lock we created.
	Unlock() error

	// intialize will check if schema_migrations and migration_locks are present or not
	// if not it create those tables
	Initialize(ctx context.Context) error

	// returns all the migrations that have run so far
	GetSchemaMigrations() ([]SchemaMigration, error)

	// adds/removes migrations from schema_migratoins table
	UpdateSchemaMigrations(migrations []SchemaMigration, migrationType string) error

	// crossCheckMigrations will match the migrations files in the file system with the files stored in schema_migrations
	// the files that are not present in schema_migrations will be executed by RunMigration Method
	CrossCheckMigrations(ctx context.Context, migrations []string, schemaMigratoins []SchemaMigration) error

	// Migrate runs a given set of migrations
	Migrate(ctx context.Context, migrationsFiles []string) error
}

// SchemaMigration defines the schema_migration table which stores all the migraitons executed
type SchemaMigration struct {
	ID            int64
	MigrationName string
	BatchNumber   int
	CreatedAt     int64
}

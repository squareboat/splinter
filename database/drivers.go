package database

import "context"

// a common driver interface that all db drivers will have to implement
type Driver interface {
	// intialize will check if schema_migrations and migration_locks are present or not
	// if not it create those tables
	Initialize(ctx context.Context) error

	// crossCheckMigrations will match the migrations files in the file system with the files stored in schema_migrations
	// the files that are not present in schema_migrations will be executed by RunMigration Method
	CrossCheckMigrations(ctx context.Context, migrationFiles []string) (map[string]string, error)

	// RunMigration runs a given set of migrations and stores the migration file names in schema_migrations
	RunMigrations(ctx context.Context, migrations map[string]string) error
}

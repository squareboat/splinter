package database

import "context"

// a common driver interface that all db drivers will have to implement
type Driver interface {
	RunMigration(ctx context.Context, migrations map[string]string) error
}

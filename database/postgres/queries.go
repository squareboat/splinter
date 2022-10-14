package postgres

import (
	"fmt"
	"strings"
	"time"

	"github.com/squareboat/splinter/database"
)

func tableExists(tableName, schemaName string) string {
	return fmt.Sprintf(`
	SELECT EXISTS (
   	SELECT FROM pg_tables
  		WHERE  schemaname = '%v'
   		AND    tablename  = '%v'
   );
	`, schemaName, tableName)
}

func createSchemaMigrations() string {
	return `
	CREATE TABLE IF NOT EXISTS schema_migrations (
						id SERIAL PRIMARY KEY,
						migration_name TEXT,
						batch_number INT,
						created_at INT8
					);
	`
}

func createMigrationLocksTable() string {
	return `
		CREATE TABLE IF NOT EXISTS migrations_lock (
			id SERIAL PRIMARY KEY,
			is_locked BOOLEAN
		);
	`
}

func getMigrations() string {
	return `
		SELECT * FROM schema_migrations ORDER BY  migration_name ASC;
	`
}

func insertMigrationLock() string {
	return `
		INSERT INTO migrations_lock (is_locked) VALUES (false);
	`
}

func updateMigrationLock(lock bool) string {
	return fmt.Sprintf(`
		UPDATE migrations_lock SET is_locked = %v
	`, lock)
}

func getLock(lockState bool) string {
	return fmt.Sprintf(`
		SELECT * FROM migrations_lock WHERE is_locked = %v;
	`, lockState)
}

func insertSchemaMigrations(migrationFiles []string, batchNumber int) string {
	var query strings.Builder
	query.WriteString("INSERT INTO schema_migrations (migration_name, batch_number, created_at)  VALUES ")

	for i := range migrationFiles {
		query.WriteString(fmt.Sprintf("( '%v' , %v, %v)", migrationFiles[i], batchNumber, time.Now().Unix()))

		if i < len(migrationFiles)-1 {
			query.WriteString(" , ")
		}
	}

	query.WriteString(";")

	return query.String()
}

func deleteLatestSchemaMigrations(filename string) string {
	return fmt.Sprintf("DELETE FROM schema_migrations WHERE migration_name = '%v'", filename)
}

func updateSchemaMigrations(migrations []database.SchemaMigration) string {
	var query strings.Builder

	query.WriteString("INSERT INTO schema_migrations (migration_name, batch_number, created_at)  VALUES ")

	for i := range migrations {
		mig := migrations[i]
		query.WriteString(fmt.Sprintf("( '%v' , %v, %v)", mig.MigrationName, mig.BatchNumber, time.Now().Unix()))

		if i < len(migrations)-1 {
			query.WriteString(" , ")
		}
	}

	return query.String()
}

func deleteFromSchemaMigrations(migrations []database.SchemaMigration) string {

	var inClause strings.Builder

	for i, mig := range migrations {
		inClause.WriteString(fmt.Sprintf(" '%v' ", mig.MigrationName))
		if i < len(migrations)-1 {
			inClause.WriteString(" , ")
		}
	}

	query := fmt.Sprintf("DELETE FROM schema_migrations WHERE migration_name IN ( %v )", inClause.String())
	return query
}

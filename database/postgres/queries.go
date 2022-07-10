package postgres

import "fmt"

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
		SELECT * FROM schema_migrations ORDER BY id DESC, migration_name ASC;
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

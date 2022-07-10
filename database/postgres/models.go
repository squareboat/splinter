package postgres

type schemaMigration struct {
	id            int64
	migrationName string
	batchNumber   int
	createdAt     int64
}

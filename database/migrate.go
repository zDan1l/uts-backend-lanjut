package database

import (
	"database/sql"
	"os"
)

// Migrate menjalankan schema.sql (CREATE TABLE IF NOT EXISTS).
func Migrate(db *sql.DB) error {
	b, err := os.ReadFile("database/schema.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(b))
	return err
}

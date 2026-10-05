package database

import (
	"database/sql"

	_ "github.com/lib/pq"
	"siakad-mini/config"
)

// Connect membuka koneksi PostgreSQL dan memastikannya hidup.
func Connect(env *config.Env) (*sql.DB, error) {
	db, err := sql.Open("postgres", env.DSN())
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

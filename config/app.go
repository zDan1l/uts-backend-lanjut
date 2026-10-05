package config

import (
	"database/sql"
	"log"
)

// App adalah wadah dependensi utama (DB + Env + Logger)
// yang dioper ke service dan route.
type App struct {
	DB  *sql.DB
	Env *Env
	Log *log.Logger
}

func New(db *sql.DB, env *Env, lg *log.Logger) *App {
	return &App{DB: db, Env: env, Log: lg}
}

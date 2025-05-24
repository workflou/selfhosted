package database

import (
	"database/sql"
	"embed"
	"log/slog"
	"os"

	_ "modernc.org/sqlite"
)

var (
	DB *sql.DB

	//go:embed schema/*.sql
	FS embed.FS
)

func init() {
	var err error

	slog.Info("[database] init", "dsn", os.Getenv("DB_DSN"))

	DB, err = sql.Open("sqlite", os.Getenv("DB_DSN"))
	if err != nil {
		panic(err)
	}

	if err = DB.Ping(); err != nil {
		panic(err)
	}

	Migrate()
}

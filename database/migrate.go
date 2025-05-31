package database

import (
	"io/fs"

	"github.com/pressly/goose/v3"
)

func Migrate() {
	goose.SetLogger(goose.NopLogger())

	schemaFS, err := fs.Sub(FS, "schema")
	if err != nil {
		panic(err)
	}

	goose.SetDialect("sqlite3")
	goose.SetBaseFS(schemaFS)

	if err = goose.Up(DB, "."); err != nil {
		panic(err)
	}
}

package app

import (
	"context"
	"template/database"
	"template/database/store"
)

var AdminCount int64 = 0

func init() {
	New()
}

func New() {
	var error error
	AdminCount, error = store.New(database.DB).CountAdmins(context.Background())

	if error != nil {
		panic(error)
	}
}

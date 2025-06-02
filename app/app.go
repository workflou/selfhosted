package app

import (
	"context"
	"os"
	"selfhosted/database"
	"selfhosted/database/store"
)

var AdminCount int64 = 0
var Name = os.Getenv("APP_NAME")

func init() {
	New()

	if Name == "" {
		Name = "selfhosted"
	}
}

func New() {
	var error error
	AdminCount, error = store.New(database.DB).CountAdmins(context.Background())

	if error != nil {
		panic(error)
	}
}

type SessionKeyType string

const SessionKey SessionKeyType = "session"

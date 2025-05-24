package test

import (
	"net/http"
	"net/http/httptest"
	"template/database"
	"template/router"

	"github.com/pressly/goose/v3"
)

type TestCase struct {
	Server *httptest.Server
	Client *http.Client
}

func NewTestCase() *TestCase {
	server := httptest.NewServer(router.New())
	client := server.Client()

	goose.Down(database.DB, ".")
	database.Migrate()

	return &TestCase{
		Server: server,
		Client: client,
	}
}

func (tc *TestCase) Close() {
	tc.Server.Close()
}

package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"template/database"
	"template/router"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/pressly/goose/v3"
)

type TestCase struct {
	Server *httptest.Server
	Client *http.Client

	T            *testing.T
	LastResponse *http.Response
	LastDocument *goquery.Document
}

func NewTestCase(t *testing.T) *TestCase {
	if os.Getenv("ENV") != "test" {
		panic("Test cases should only be run in the test environment")
	}

	server := httptest.NewServer(router.New())
	client := server.Client()

	goose.Down(database.DB, ".")
	database.Migrate()

	return &TestCase{
		T:      t,
		Server: server,
		Client: client,
	}
}

func (tc *TestCase) Close() {
	tc.Server.Close()
}

func (tc *TestCase) Get(path string) (*http.Response, error) {
	res, err := tc.Client.Get(tc.Server.URL + path)
	if err != nil {
		return nil, err
	}

	tc.LastResponse = res
	tc.LastDocument, _ = goquery.NewDocumentFromReader(res.Body)

	return res, nil
}

func (tc *TestCase) AssertRedirect(statusCode int, location string) {
	if tc.LastResponse == nil {
		tc.T.Fatal("No response available to assert redirect")
	}

	if tc.LastResponse.Request.Response.StatusCode != http.StatusSeeOther {
		tc.T.Fatalf("Expected status code %d, got %d", http.StatusSeeOther, tc.LastResponse.Request.Response.StatusCode)
	}

	if tc.LastResponse.Request.Response.Header.Get("Location") != "/setup" {
		tc.T.Fatalf("Expected redirect to %s, got %s", "/setup", tc.LastResponse.Request.Response.Header.Get("Location"))
	}
}

func (tc *TestCase) AssertElementVisible(selector string) {
	if tc.LastDocument == nil {
		tc.T.Fatal("No response available to assert element visibility")
	}

	selection := tc.LastDocument.Find(selector)
	if selection.Length() == 0 {
		tc.T.Fatalf(
			"Expected element with selector '%s' to be visible, but it was not found. The output was:\n%s",
			selector,
			tc.LastDocument.Text(),
		)
	}
}

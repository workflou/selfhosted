package test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"selfhosted/router"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/pressly/goose/v3"
	"golang.org/x/crypto/bcrypt"
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

	app.New()

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

func (tc *TestCase) Post(path string, body url.Values) (*http.Response, error) {
	res, err := tc.Client.PostForm(tc.Server.URL+path, body)
	if err != nil {
		return nil, err
	}

	tc.LastResponse = res
	tc.LastDocument, _ = goquery.NewDocumentFromReader(res.Body)

	return res, nil
}

func (tc *TestCase) SetupAdmin() {
	store.New(database.DB).CreateAdmin(tc.T.Context(), store.CreateAdminParams{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: "$2a$10$EIX/3Z1z5Q8b1Y4e5f6e9O0j7k5h5F5y5F5y5F5y5F5y5F5y5F5y", // password123
	})
	app.AdminCount = 1
}

func (tc *TestCase) CreateUser(name, email, password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		tc.T.Fatalf("Failed to hash password: %v", err)
	}

	err = store.New(database.DB).CreateUser(tc.T.Context(), store.CreateUserParams{
		Name:     name,
		Email:    strings.ToLower(email),
		Password: string(hash),
	})

	if err != nil {
		tc.T.Fatalf("Failed to create user: %v", err)
	}
}

func (tc *TestCase) AssertRedirect(statusCode int, location string) {
	if tc.LastResponse == nil {
		tc.T.Fatal("No response available to assert redirect")
	}

	if tc.LastResponse.Request == nil || tc.LastResponse.Request.Response == nil {
		tc.T.Fatal("No request or response available to assert redirect")
	}

	if tc.LastResponse.Request.Response.StatusCode != http.StatusSeeOther {
		tc.T.Fatalf("Expected status code %d, got %d", http.StatusSeeOther, tc.LastResponse.Request.Response.StatusCode)
	}

	if tc.LastResponse.Request.Response.Header.Get("Location") != location {
		tc.T.Fatalf("Expected redirect to %s, got %s", location, tc.LastResponse.Request.Response.Header.Get("Location"))
	}
}

func (tc *TestCase) AssertStatus(expectedStatus int) {
	if tc.LastResponse == nil {
		tc.T.Fatal("No response available to assert status")
	}

	if tc.LastResponse.StatusCode != expectedStatus {
		tc.T.Fatalf("Expected status code %d, got %d", expectedStatus, tc.LastResponse.StatusCode)
	}
}

func (tc *TestCase) AssertHeader(header, expectedValue string) {
	if tc.LastResponse == nil {
		tc.T.Fatal("No response available to assert header")
	}

	if tc.LastResponse.Request != nil && tc.LastResponse.Request.Response != nil {
		actualValue := tc.LastResponse.Request.Response.Header.Get(header)
		if actualValue == expectedValue {
			return
		}
	}

	actualValue := tc.LastResponse.Header.Get(header)
	if actualValue != expectedValue {
		tc.T.Fatalf("Expected header '%s' to be '%s', got '%s'", header, expectedValue, actualValue)
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

func (tc *TestCase) AssertDatabaseCount(table string, expectedCount int) {
	query := "SELECT COUNT(*) FROM " + table
	row := database.DB.QueryRow(query)
	var count int
	err := row.Scan(&count)
	if err != nil {
		tc.T.Fatalf("Failed to query database for count of table '%s': %v", table, err)
	}
	if count != expectedCount {
		tc.T.Fatalf("Expected %d rows in table '%s', got %d", expectedCount, table, count)
	}
}

func (tc *TestCase) AssertDatabaseMissing(table string, conditions map[string]interface{}) {
	query := "SELECT COUNT(*) FROM " + table + " WHERE "
	args := make([]interface{}, 0, len(conditions))
	clauses := make([]string, 0, len(conditions))

	for column, value := range conditions {
		clauses = append(clauses, column+" = ?")
		args = append(args, value)
	}

	query += strings.Join(clauses, " AND ")
	row := database.DB.QueryRow(query, args...)
	var count int
	err := row.Scan(&count)
	if err != nil {
		tc.T.Fatalf("Failed to query database for missing condition in table '%s': %v", table, err)
	}
	if count > 0 {
		tc.T.Fatalf("Expected no rows in table '%s' matching conditions %v, but found %d", table, conditions, count)
	}
}

func (tc *TestCase) AssertDatabaseHas(table string, conditions map[string]interface{}) {
	query := "SELECT COUNT(*) FROM " + table + " WHERE "
	args := make([]interface{}, 0, len(conditions))
	clauses := make([]string, 0, len(conditions))

	for column, value := range conditions {
		clauses = append(clauses, column+" = ?")
		args = append(args, value)
	}

	query += strings.Join(clauses, " AND ")
	row := database.DB.QueryRow(query, args...)
	var count int
	err := row.Scan(&count)
	if err != nil {
		tc.T.Fatalf("Failed to query database for existence in table '%s': %v", table, err)
	}
	if count == 0 {
		tc.T.Fatalf("Expected at least one row in table '%s' matching conditions %v, but found none", table, conditions)
	}
}

func (tc *TestCase) AssertCookieSet(name string) {
	if tc.LastResponse == nil {
		tc.T.Fatal("No response available to assert cookie")
	}

	cookies := tc.LastResponse.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			return
		}
	}

	tc.T.Fatalf("Expected cookie '%s' to be set, but it was not found in the response", name)
}

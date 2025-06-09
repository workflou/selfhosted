package test

import (
	"io"
	"mime/multipart"
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
	UserCookie   *http.Cookie
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

	os.RemoveAll("./uploads")
}

func (tc *TestCase) Get(path string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", tc.Server.URL+path, nil)
	if tc.UserCookie != nil {
		req.AddCookie(tc.UserCookie)
	}

	res, err := tc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	tc.LastResponse = res
	tc.LastDocument, _ = goquery.NewDocumentFromReader(res.Body)

	return res, nil
}

func (tc *TestCase) Post(path string, body url.Values) (*http.Response, error) {
	req, _ := http.NewRequest("POST", tc.Server.URL+path, strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if tc.UserCookie != nil {
		req.AddCookie(tc.UserCookie)
	}

	res, err := tc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	tc.LastResponse = res
	tc.LastDocument, _ = goquery.NewDocumentFromReader(res.Body)

	return res, nil
}

func (tc *TestCase) PostMultipart(path string, reader io.Reader, writer *multipart.Writer) (*http.Response, error) {
	req, _ := http.NewRequest("POST", tc.Server.URL+path, reader)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if tc.UserCookie != nil {
		req.AddCookie(tc.UserCookie)
	}

	res, err := tc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	tc.LastResponse = res
	tc.LastDocument, _ = goquery.NewDocumentFromReader(res.Body)

	return res, nil
}

func (tc *TestCase) SetupAdmin() {
	hash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		tc.T.Fatalf("Failed to hash password: %v", err)
	}

	store.New(database.DB).CreateAdmin(tc.T.Context(), store.CreateAdminParams{
		Name:     "Admin",
		Email:    "admin@example.com",
		Password: string(hash),
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

func (tc *TestCase) LogIn(email, password string) {
	res, _ := tc.Post("/login", url.Values{
		"email":    {email},
		"password": {password},
	})

	if res.StatusCode != http.StatusOK {
		tc.T.Fatalf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	for _, cookie := range res.Cookies() {
		if cookie.Name == "session" {
			tc.UserCookie = cookie
		}
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

func (tc *TestCase) AssertNoRedirect() {
	if tc.LastResponse == nil {
		tc.T.Fatalf("No response available to assert no redirect")
	}

	if tc.LastResponse.Request != nil && tc.LastResponse.Request.Response != nil {
		tc.T.Fatalf("Expected no redirect, but got status code %d", tc.LastResponse.Request.Response.StatusCode)
	}
}

func (tc *TestCase) AssertStatus(expectedStatus int) {
	if tc.LastResponse == nil {
		tc.T.Fatal("No response available to assert status")
	}

	if tc.LastResponse.StatusCode != expectedStatus {
		body, _ := io.ReadAll(tc.LastResponse.Body)
		tc.T.Fatalf("Expected status code %d, got %d. Response body: %s", expectedStatus, tc.LastResponse.StatusCode, body)
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

func (tc *TestCase) AssertCookieMissing(name string) {
	if tc.LastResponse == nil {
		tc.T.Fatal("No response available to assert cookie")
	}

	cookies := tc.LastResponse.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == name {
			tc.T.Fatalf("Expected cookie '%s' to be missing, but it was found in the response", name)
		}
	}
}

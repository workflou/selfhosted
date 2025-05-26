package test

import (
	"context"
	"net/http"
	"net/url"
	"template/app"
	"template/database"
	"template/database/store"
	"testing"
)

func TestSetup(t *testing.T) {
	t.Run("redirect to setup when no admin account exists", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/")
		tc.AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("setup form is displayed", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/setup")

		tc.AssertElementVisible("form[hx-post='/setup']")
		tc.AssertElementVisible("input[name='name']")
		tc.AssertElementVisible("input[name='email']")
		tc.AssertElementVisible("input[name='password']")
		tc.AssertElementVisible("button[type='submit']")
	})

	t.Run("setup form validation", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Post("/setup", url.Values{})
		tc.AssertStatus(http.StatusBadRequest)
	})

	t.Run("admin account can be created", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		formData := url.Values{
			"name":     {"Admin"},
			"email":    {"admin@example.com"},
			"password": {"password123"},
		}

		tc.Post("/setup", formData)

		tc.AssertStatus(http.StatusOK)
		tc.AssertHeader("HX-Redirect", "/")
		tc.AssertDatabaseCount("users", 1)
	})

	t.Run("password is not stored in plaintext", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		formData := url.Values{
			"name":     {"Admin"},
			"email":    {"admin@example.com"},
			"password": {"password123"},
		}

		tc.Post("/setup", formData)

		tc.AssertDatabaseMissing("users", map[string]interface{}{
			"password": "password123",
		})
	})

	t.Run("user already exists", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		store.New(database.DB).CreateAdmin(context.Background(), store.CreateAdminParams{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: "$2a$10$EIX/3Z1z5Q8b1Y4e5f6e9O0j7k5h5F5y5F5y5F5y5F5y5F5y5F5y",
		})

		formData := url.Values{
			"name":     {"Admin"},
			"email":    {"admin2@example.com"},
			"password": {"password123"},
		}

		tc.Post("/setup", formData)
		tc.AssertStatus(http.StatusConflict)
	})

	t.Run("email is lowercased", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		formData := url.Values{
			"name":     {"Admin"},
			"email":    {"AdMiN@eXaMpLe.CoM"},
			"password": {"password123"},
		}

		tc.Post("/setup", formData)

		tc.AssertDatabaseHas("users", map[string]interface{}{
			"email": "admin@example.com",
		})

		tc.AssertDatabaseMissing("users", map[string]interface{}{
			"email": "AdMiN@eXaMpLe.CoM",
		})
	})

	t.Run("admin count is updated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		if app.AdminCount != 0 {
			t.Fatalf("Expected app.AdminCount to be 0, got %d", app.AdminCount)
		}

		formData := url.Values{
			"name":     {"Admin"},
			"email":    {"admin2@example.com"},
			"password": {"password123"},
		}

		tc.Post("/setup", formData)

		if app.AdminCount != 1 {
			t.Fatalf("Expected app.AdminCount to be 1, got %d", app.AdminCount)
		}
	})

	t.Run("setup form redirects to home after successful setup", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		store.New(database.DB).CreateAdmin(context.Background(), store.CreateAdminParams{
			Name:     "Admin",
			Email:    "admin@example.com",
			Password: "$2a$10$EIX/3Z1z5Q8b1Y4e5f6e9O0j7k5h5F5y5F5y5F5y5F5y5F5y5F5y",
		})

		app.AdminCount = 1

		tc.Get("/setup")
		tc.AssertRedirect(http.StatusSeeOther, "/login")
	})
}

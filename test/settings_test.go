package test

import (
	"context"
	"net/http"
	"net/url"
	"selfhosted/database"
	"selfhosted/database/store"
	"testing"
)

func TestSettings(t *testing.T) {
	t.Run("guests can't access settings page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.Get("/settings")
		tc.AssertRedirect(http.StatusFound, "/login")
	})

	t.Run("settings page is displayed for logged-in users", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.LogIn("admin@example.com", "password123")
		tc.Get("/settings")
		tc.AssertStatus(http.StatusOK)
		tc.AssertElementVisible("form[hx-post='/settings']")
		tc.AssertElementVisible("input[name='name']")
		tc.AssertElementVisible("input[name='avatar']")
		tc.AssertElementVisible("button[type='submit']")
	})

	t.Run("name can be updated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.LogIn("admin@example.com", "password123")

		formData := url.Values{
			"name": {"New Admin Name"},
		}

		tc.Post("/settings", formData)
		tc.AssertStatus(http.StatusOK)
		tc.AssertHeader("HX-Reswap", "none")

		tc.AssertDatabaseCount("users", 1)
		tc.AssertDatabaseHas("users", map[string]any{
			"name": "New Admin Name",
		})
	})

	t.Run("avatar can be updated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.LogIn("admin@example.com", "password123")

		formData := url.Values{
			"name":   {"Admin with Avatar"},
			"avatar": {"testdata/avatar.jpeg"},
		}

		tc.PostMultipart("/settings", formData)

		tc.AssertStatus(http.StatusOK)
		tc.AssertHeader("HX-Reswap", "none")
		tc.AssertDatabaseCount("users", 1)

		user, err := store.New(database.DB).GetUserByEmail(context.Background(), "admin@example.com")
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		if user.Avatar.String == "" {
			t.Error("expected avatar to be set, but it was empty")
		}
	})
}

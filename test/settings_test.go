package test

import (
	"net/http"
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
}

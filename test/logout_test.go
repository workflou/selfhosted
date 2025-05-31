package test

import (
	"net/http"
	"testing"
)

func TestLogout(t *testing.T) {
	t.Run("guests cannot logout", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.Get("/logout")
		tc.AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("logged in users can logout", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.LogIn("admin@example.com", "password123")

		tc.Get("/logout")
		tc.AssertRedirect(http.StatusSeeOther, "/login")
		tc.AssertCookieMissing("session")

		tc.AssertDatabaseCount("sessions", 0)

		tc.Get("/")
		tc.AssertRedirect(http.StatusSeeOther, "/login")
	})
}

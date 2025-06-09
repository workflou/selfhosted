package test

import (
	"net/http"
	"testing"
)

func TestUserMenu(t *testing.T) {
	t.Run("guests cannot access user modal", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.Get("/user")
		tc.AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("users can access user modal", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.LogIn("admin@example.com", "password123")
		tc.Get("/user")
		tc.AssertStatus(http.StatusOK)
		tc.AssertNoRedirect()
		tc.AssertElementVisible("nav")
	})
}

package test

import "testing"

func TestTeam(t *testing.T) {
	t.Run("user can see a current team", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.LogIn("admin@example.com", "password123")

		tc.Get("/")
		tc.AssertSee("Admin Team")
	})
}

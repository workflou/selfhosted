package test

import (
	"net/url"
	"selfhosted/cmd"
	"testing"
)

func TestAdminPassword(t *testing.T) {
	t.Run("admin password can be changed", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		cmd.ChangeAdminPassword("admin@example.com", "newpassword123")

		tc.Post("/login", url.Values{
			"email":    {"admin@example.com"},
			"password": {"password123"},
		})
		tc.AssertStatus(400)

		tc.Post("/login", url.Values{
			"email":    {"admin@example.com"},
			"password": {"newpassword123"},
		})
		tc.AssertStatus(200)
	})
}

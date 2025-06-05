package test

import (
	"net/http"
	"testing"
)

func TestSettings(t *testing.T) {
	t.Run("guests can't access settings page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/settings")
		tc.AssertRedirect(http.StatusFound, "/login")
	})
}

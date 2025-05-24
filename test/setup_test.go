package test

import (
	"net/http"
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
}

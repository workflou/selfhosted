package test

import (
	"net/http"
	"testing"
)

func TestSetupRedirect(t *testing.T) {
	tc := NewTestCase(t)
	defer tc.Close()

	tc.Get("/")
	tc.AssertRedirect(http.StatusSeeOther, "/setup")
}

func TestSetupForm(t *testing.T) {
	tc := NewTestCase(t)
	defer tc.Close()

	tc.Get("/setup")

	tc.AssertElementVisible("form[hx-post='/setup']")
	tc.AssertElementVisible("input[name='name']")
	tc.AssertElementVisible("input[name='email']")
	tc.AssertElementVisible("input[name='password']")
	tc.AssertElementVisible("button[type='submit']")
}

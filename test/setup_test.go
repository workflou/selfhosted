package test

import (
	"net/http"
	"testing"
)

func TestSetupRedirect(t *testing.T) {
	tc := NewTestCase()
	defer tc.Close()

	res, _ := tc.Client.Get(tc.Server.URL + "/")

	if res.Request.Response.StatusCode != http.StatusSeeOther {
		t.Fatalf("Expected status code %d, got %d", http.StatusSeeOther, res.Request.Response.StatusCode)
	}

	if res.Request.Response.Header.Get("Location") != "/setup" {
		t.Fatalf("Expected redirect to %s, got %s", "/setup", res.Request.Response.Header.Get("Location"))
	}
}

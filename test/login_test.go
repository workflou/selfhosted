package test

import (
	"context"
	"net/http"
	"net/url"
	"selfhosted/database"
	"selfhosted/database/store"
	"testing"
	"time"
)

func TestLogin(t *testing.T) {
	t.Run("if no admin account exists, redirect to setup", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.Get("/login")
		tc.AssertRedirect(http.StatusSeeOther, "/setup")
	})

	t.Run("guests are redirected to login page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()
		tc.SetupAdmin()

		tc.Get("/")
		tc.AssertRedirect(http.StatusSeeOther, "/login")
	})

	t.Run("login page is displayed", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.Get("/login")
		tc.AssertStatus(http.StatusOK)

		tc.AssertElementVisible("form[hx-post='/login']")
		tc.AssertElementVisible("input[name='email']")
		tc.AssertElementVisible("input[name='password']")
		tc.AssertElementVisible("button[type='submit']")
	})

	t.Run("login form validation", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.Post("/login", nil)
		tc.AssertStatus(http.StatusBadRequest)

		tc.Post("/login", url.Values{
			"email":    {"invalid"},
			"password": {"password123"},
		})
		tc.AssertStatus(http.StatusBadRequest)
	})

	t.Run("user has to exist", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		tc.Post("/login", url.Values{
			"email":    {"notfound@example.com"},
			"password": {"password123"},
		})

		tc.AssertStatus(http.StatusBadRequest)
	})

	t.Run("password has to match", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.CreateUser("User", "user@example.com", "password123")

		tc.Post("/login", url.Values{
			"email":    {"user@example.com"},
			"password": {"wrongpassword"},
		})

		tc.AssertStatus(http.StatusBadRequest)
	})

	t.Run("user can log in", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.CreateUser("User", "user@example.com", "password123")

		res, _ := tc.Post("/login", url.Values{
			"email":    {"user@example.com"},
			"password": {"password123"},
		})

		tc.AssertStatus(http.StatusOK)
		tc.AssertHeader("HX-Push-Url", "/")
		tc.AssertHeader("HX-Trigger", `{"redirect": "/"}`)
		tc.AssertCookieSet("session")

		tc.AssertDatabaseCount("sessions", 1)

		var responseCookie *http.Cookie

		for _, cookie := range res.Cookies() {
			if cookie.Name == "session" {
				if cookie.HttpOnly == false {
					t.Errorf("session cookie should be HttpOnly")
				}

				if cookie.Value == "" {
					t.Errorf("session cookie should have a value")
				}

				responseCookie = cookie

				s, err := store.New(database.DB).GetSessionByUuid(context.Background(), cookie.Value)
				if err != nil || s.ID == 0 {
					t.Errorf("failed to get session by ID: %v", err)
				}

				if cookie.Expires.Round(time.Minute).Equal(s.ExpiresAt.Round(time.Minute)) == false {
					t.Errorf("session cookie expiration does not match session expiration. Expected: %v, got: %v", s.ExpiresAt.Round(time.Minute), cookie.Expires.Round(time.Minute))
				}

				if cookie.SameSite != http.SameSiteLaxMode {
					t.Errorf("session cookie SameSite should be Lax, got %v", cookie.SameSite)
				}
			}
		}

		req, _ := http.NewRequest(http.MethodGet, tc.Server.URL+"/", nil)
		req.AddCookie(responseCookie)

		res, err := tc.Client.Do(req)
		if err != nil {
			t.Fatalf("failed to make request: %v", err)
		}

		if res.StatusCode != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
		}

		if res.Request.Response != nil {
			t.Fatalf("expected no redirect, got %s to %s", res.Request.Response.Request.URL, res.Request.Response.Header.Get("Location"))
		}
	})

	t.Run("login is rate limited", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()

		for i := 0; i < 5; i++ {
			tc.Post("/login", url.Values{
				"email":    {"invalid@example.com"},
				"password": {"wrongpassword"},
			})
			tc.AssertStatus(http.StatusBadRequest)
		}

		tc.Post("/login", url.Values{
			"email":    {"invalid@example.com"},
			"password": {"wrongpassword"},
		})
		tc.AssertStatus(http.StatusTooManyRequests)
		tc.AssertHeader("Retry-After", "60")
		tc.AssertHeader("X-RateLimit-Limit", "5")
		tc.AssertHeader("X-RateLimit-Remaining", "0")
	})

	t.Run("authenticated users cannot access login page", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.CreateUser("User", "user@example.com", "password123")

		tc.LogIn("user@example.com", "password123")

		tc.Get("/login")
		tc.AssertRedirect(http.StatusSeeOther, "/")
	})
}

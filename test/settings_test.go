package test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"selfhosted/database"
	"selfhosted/database/store"
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
		tc.AssertElementVisible("input[name='name']")
		tc.AssertElementVisible("input[name='avatar']")
	})

	t.Run("name can be updated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.LogIn("admin@example.com", "password123")

		formData := url.Values{
			"name": {"New Admin Name"},
		}

		tc.Post("/settings/name", formData)
		tc.AssertStatus(http.StatusOK)
		tc.AssertHeader("HX-Reswap", "none")

		tc.AssertDatabaseCount("users", 1)
		tc.AssertDatabaseHas("users", map[string]any{
			"name": "New Admin Name",
		})
	})

	t.Run("avatar can be updated", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.LogIn("admin@example.com", "password123")

		f, err := os.Open("./../testdata/avatar.jpeg")
		if err != nil {
			t.Fatalf("failed to open avatar file: %v", err)
		}
		defer f.Close()

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		avatar, err := writer.CreateFormFile("avatar", "avatar.jpeg")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		_, err = io.Copy(avatar, f)
		if err != nil {
			t.Fatalf("failed to copy avatar file: %v", err)
		}

		err = writer.Close()
		if err != nil {
			t.Fatalf("failed to close multipart writer: %v", err)
		}

		tc.PostMultipart("/settings/avatar", body, writer)

		tc.AssertStatus(http.StatusOK)
		tc.AssertHeader("HX-Reswap", "none")
		tc.AssertDatabaseCount("users", 1)

		user, err := store.New(database.DB).GetUserByEmail(context.Background(), "admin@example.com")
		if err != nil {
			t.Fatalf("failed to get user: %v", err)
		}

		if user.Avatar.String == "" {
			t.Error("expected avatar to be set, but it was empty")
		}
	})

	t.Run("invalid avatar upload is rejected", func(t *testing.T) {
		tc := NewTestCase(t)
		defer tc.Close()

		tc.SetupAdmin()
		tc.LogIn("admin@example.com", "password123")

		body := new(bytes.Buffer)
		writer := multipart.NewWriter(body)

		nonImageFile, err := writer.CreateFormFile("avatar", "document.txt")
		if err != nil {
			t.Fatalf("failed to create form file: %v", err)
		}
		_, err = nonImageFile.Write([]byte("This is a text file, not an image."))
		if err != nil {
			t.Fatalf("failed to write to form file: %v", err)
		}
		err = writer.Close()
		if err != nil {
			t.Fatalf("failed to close multipart writer: %v", err)
		}
		tc.PostMultipart("/settings/avatar", body, writer)
		tc.AssertStatus(http.StatusBadRequest)
	})
}

package handler

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"selfhosted/html"
	"selfhosted/toast"

	"github.com/google/uuid"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	html.SettingsPage().Render(r.Context(), w)
}

func SettingsNameForm(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContext(r.Context())
	if user == nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusUnauthorized)
		toast.Error("Unauthorized", "You must be logged in to update your profile.").Send(r.Context(), w)
		return
	}

	r.ParseForm()

	name := r.FormValue("name")
	if name == "" {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusBadRequest)
		toast.Error("Invalid request", "Name is required").Send(r.Context(), w)
		return
	}

	err := store.New(database.DB).UpdateUserName(r.Context(), store.UpdateUserNameParams{
		ID:   user.ID,
		Name: name,
	})
	if err != nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusInternalServerError)
		toast.Error("Update failed", "An error occurred while updating your profile. Please try again later.").Send(r.Context(), w)
		return
	}

	user.Name = name
	w.Header().Set("HX-Reswap", "none")
	html.UserName(user.Name).Render(r.Context(), w)
	toast.Success("Name updated", "Your name has been successfully updated.").Send(r.Context(), w)
}

func SettingsAvatarForm(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContext(r.Context())
	if user == nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusUnauthorized)
		toast.Error("Unauthorized", "You must be logged in to update your profile.").Send(r.Context(), w)
		return
	}

	r.ParseMultipartForm(10 << 20)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	path, err := uploadFile(r, "avatar", "./uploads/avatars")
	if err != nil && err != http.ErrMissingFile {
		slog.Error("File upload error", "error", err)
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusInternalServerError)
		toast.Error("Upload failed", err.Error()).Send(r.Context(), w)
		return
	}

	if err != nil {
		slog.Error("File upload error", "error", err)
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusInternalServerError)
		toast.Error("Upload failed", err.Error()).Send(r.Context(), w)
		return
	}

	if user.Avatar.Valid {
		err = os.Remove(filepath.Join(".", user.Avatar.String))
		if err != nil {
			slog.Error("Failed to remove old avatar", "error", err)
		}
	}

	err = store.New(database.DB).UpdateUserAvatar(r.Context(), store.UpdateUserAvatarParams{
		ID: user.ID,
		Avatar: sql.NullString{
			String: "/" + path,
			Valid:  true,
		},
	})
	if err != nil {
		slog.Error("Failed to update user avatar", "error", err)
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusInternalServerError)
		toast.Error("Update failed", "An error occurred while updating your profile. Please try again later.").Send(r.Context(), w)
		return
	}
	user.Avatar = sql.NullString{
		String: "/" + path,
		Valid:  true,
	}
	html.UserAvatar(user.Avatar.String, user.Name).Render(r.Context(), w)
	toast.Success("Avatar updated", "Your avatar has been successfully updated.").Send(r.Context(), w)
	return
}

func uploadFile(r *http.Request, fieldName string, uploadsPath string) (string, error) {
	file, fileHeader, err := r.FormFile(fieldName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = os.MkdirAll(uploadsPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	filename := uuid.New().String() + filepath.Ext(fileHeader.Filename)
	path := filepath.Join(uploadsPath, filename)

	dst, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("failed to create upload file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to copy uploaded file: %w", err)
	}

	return path, nil
}

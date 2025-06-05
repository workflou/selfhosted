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

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	user := app.GetUserFromContext(r.Context())
	if user == nil {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusUnauthorized)
		toast.Error("Unauthorized", "You must be logged in to update your profile.").Send(r.Context(), w)
		return
	}

	path, err := uploadFile(r, "avatar", "./uploads/avatars")
	if err != nil && err != http.ErrMissingFile {
		slog.Error("File upload error", "error", err)
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusInternalServerError)
		toast.Error("Upload failed", err.Error()).Send(r.Context(), w)
		return
	}

	if err == nil {
		slog.Info("User avatar updated", "user_id", user.ID, "avatar_path", path)
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
	}

	if r.FormValue("name") != "" {
		err = store.New(database.DB).UpdateUserName(r.Context(), store.UpdateUserNameParams{
			ID:   user.ID,
			Name: r.FormValue("name"),
		})
		if err != nil {
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusInternalServerError)
			toast.Error("Update failed", "An error occurred while updating your profile. Please try again later.").Send(r.Context(), w)
			return
		}
		user.Name = r.FormValue("name")
	}

	w.Header().Set("HX-Reswap", "none")
	toast.Success("Profile updated", "Your profile has been successfully updated.").Send(r.Context(), w)
	html.UserName(user.Name).Render(r.Context(), w)
	html.UserAvatar(user.Avatar.String, user.Name).Render(r.Context(), w)
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

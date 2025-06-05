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

	file, fileHeader, err := r.FormFile("avatar")
	if err == nil {
		defer file.Close()
		filename := fileHeader.Filename

		uploads := "./uploads/avatars"
		err = os.MkdirAll(uploads, os.ModePerm)
		if err != nil {
			slog.Error("Failed to create upload directory", "error", err)
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusInternalServerError)
			toast.Error("Upload failed", "An error occurred while creating the upload directory.").Send(r.Context(), w)
			return
		}

		uuid := uuid.New().String()
		path := filepath.Join(uploads, fmt.Sprintf("%s%s", uuid, filepath.Ext(filename)))

		dst, err := os.Create(path)
		if err != nil {
			slog.Error("Failed to create upload file", "error", err)
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusInternalServerError)
			toast.Error("Upload failed", "An error occurred while saving the uploaded file.").Send(r.Context(), w)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			slog.Error("Failed to copy uploaded file", "error", err)
			w.Header().Set("HX-Reswap", "none")
			w.WriteHeader(http.StatusInternalServerError)
			toast.Error("Upload failed", "An error occurred while copying the uploaded file.").Send(r.Context(), w)
			return
		}

		slog.Info("User avatar updated", "user_id", user.ID, "avatar_path", dst.Name())
		err = store.New(database.DB).UpdateUserAvatar(r.Context(), store.UpdateUserAvatarParams{
			ID: user.ID,
			Avatar: sql.NullString{
				String: "/" + dst.Name(),
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
			String: "/" + dst.Name(),
			Valid:  true,
		}
	}

	if r.FormValue("name") == "" {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusBadRequest)
		toast.Error("Invalid input", "Please provide a valid name.").Send(r.Context(), w)
		return
	}

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

	w.Header().Set("HX-Reswap", "none")
	toast.Success("Profile updated", "Your profile has been successfully updated.").Send(r.Context(), w)
	html.UserName(user.Name).Render(r.Context(), w)
	html.UserAvatar(user.Avatar.String, user.Name).Render(r.Context(), w)
}

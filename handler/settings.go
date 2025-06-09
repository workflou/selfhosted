package handler

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"selfhosted/html"
	"selfhosted/pkg/file"
	"selfhosted/toast"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	html.SettingsPage().Render(r.Context(), w)
}

func SettingsNameForm(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContext(r.Context())

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
	html.UserName(user.Name, true).Render(r.Context(), w)
	toast.Success("Name updated", "Your name has been successfully updated.").Send(r.Context(), w)
}

func SettingsAvatarForm(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContext(r.Context())

	r.ParseMultipartForm(10 << 20)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if !file.ValidateImageFromRequest(r, "avatar") {
		w.Header().Set("HX-Reswap", "none")
		w.WriteHeader(http.StatusBadRequest)
		toast.Error("Invalid file type", "Only image files are allowed.").Send(r.Context(), w)
		return
	}

	path, err := file.UploadFromRequest(r, "avatar", "./uploads/avatars")
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
	w.Header().Set("HX-Reswap", "none")
	html.UserAvatar(user.Avatar.String, user.Name, true).Render(r.Context(), w)
	toast.Success("Avatar updated", "Your avatar has been successfully updated.").Send(r.Context(), w)
	return
}

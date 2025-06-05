package handler

import (
	"net/http"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"selfhosted/html"
	"selfhosted/toast"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	html.SettingsPage().Render(r.Context(), w)
}

func SettingsForm(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	if r.FormValue("avatar") != "" {
		file, _, err := r.FormFile("avatar")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			toast.Error("Invalid file", "Please upload a valid image file.").Send(r.Context(), w)
			html.SettingsForm().Render(r.Context(), w)
			return
		}
		defer file.Close()

	}

	if r.FormValue("name") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toast.Error("Invalid input", "Please provide a valid name.").Send(r.Context(), w)
		html.SettingsForm().Render(r.Context(), w)
		return
	}

	user := app.GetUserFromContext(r.Context())
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		toast.Error("Unauthorized", "You must be logged in to update your profile.").Send(r.Context(), w)
		html.SettingsForm().Render(r.Context(), w)
		return
	}

	err := store.New(database.DB).UpdateUserName(r.Context(), store.UpdateUserNameParams{
		ID:   user.ID,
		Name: r.FormValue("name"),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		toast.Error("Update failed", "An error occurred while updating your profile. Please try again later.").Send(r.Context(), w)
		html.SettingsForm().Render(r.Context(), w)
		return
	}

	w.Header().Set("HX-Reswap", "none")
	toast.Success("Profile updated", "Your profile has been successfully updated.").Send(r.Context(), w)
	html.UserName(r.FormValue("name")).Render(r.Context(), w)
	html.SettingsForm().Render(r.Context(), w)
}

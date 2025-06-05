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
	r.ParseForm()

	if r.FormValue("name") == "" {
		w.WriteHeader(http.StatusBadRequest)
		html.SettingsForm().Render(r.Context(), w)
		return
	}

	user := app.GetUserFromContext(r.Context())
	if user == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	err := store.New(database.DB).UpdateUserName(r.Context(), store.UpdateUserNameParams{
		ID:   user.ID,
		Name: r.FormValue("name"),
	})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Reswap", "none")
	toast.Success("Profile updated", "Your profile has been successfully updated.").Send(r.Context(), w)
	html.SettingsForm().Render(r.Context(), w)
}

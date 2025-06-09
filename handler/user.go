package handler

import (
	"net/http"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"selfhosted/html"
)

func UserModal(w http.ResponseWriter, r *http.Request) {
	user := app.GetUserFromContext(r.Context())
	if user == nil {
		w.Header().Set("HX-Location", "/login")
		return
	}

	teams, err := store.New(database.DB).GetUserTeams(r.Context(), user.ID)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	html.UserModal(html.UserModalProps{
		Teams: teams,
	}).Render(r.Context(), w)
}

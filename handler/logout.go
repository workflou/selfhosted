package handler

import (
	"net/http"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	sess := app.GetSessionFromContext(r.Context())
	if sess == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	store.New(database.DB).DeleteSession(r.Context(), sess.Uuid)

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

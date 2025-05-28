package handler

import (
	"net/http"
	"net/mail"
	"template/database"
	"template/database/store"
	"time"

	"github.com/go-chi/httprate"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var loginRateLimiter = httprate.NewRateLimiter(5, time.Minute)

func LoginForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("email") == "" || r.FormValue("password") == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")

	if loginRateLimiter.RespondOnLimit(w, r, email) {
		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	u, err := store.New(database.DB).GetUserByEmail(r.Context(), email)
	if err != nil || u.ID == 0 {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.FormValue("password")))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	uuid := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour * 30).UTC()

	err = store.New(database.DB).CreateSession(r.Context(), store.CreateSessionParams{
		Uuid:      uuid,
		UserID:    u.ID,
		ExpiresAt: expiresAt,
	})

	w.Header().Set("HX-Redirect", "/")

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		HttpOnly: true,
		Value:    uuid,
		Expires:  expiresAt,
		SameSite: http.SameSiteLaxMode,
	})
}

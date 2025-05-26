package handler

import (
	"net/http"
	"net/mail"
	"strings"
	"template/database"
	"template/database/store"

	"golang.org/x/crypto/bcrypt"
)

func SetupForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("name") == "" || r.FormValue("email") == "" || r.FormValue("password") == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_, err := mail.ParseAddress(r.FormValue("email"))
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	count, err := store.New(database.DB).CountAdmins(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	store.New(database.DB).CreateAdmin(r.Context(), store.CreateAdminParams{
		Name:     r.FormValue("name"),
		Email:    strings.ToLower(r.FormValue("email")),
		Password: string(hash),
	})

	w.Header().Set("HX-Redirect", "/")
}

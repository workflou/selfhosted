package handler

import (
	"net/http"
	"net/mail"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"selfhosted/html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func SetupPage(w http.ResponseWriter, r *http.Request) {
	count, err := store.New(database.DB).CountAdmins(r.Context())
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if count > 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	html.SetupPage().Render(r.Context(), w)
}

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

	err = store.New(database.DB).CreateAdmin(r.Context(), store.CreateAdminParams{
		Name:     r.FormValue("name"),
		Email:    strings.ToLower(r.FormValue("email")),
		Password: string(hash),
	})

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	app.AdminCount = 1

	w.Header().Set("HX-Redirect", "/")
}

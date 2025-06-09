package handler

import (
	"log/slog"
	"net/http"
	"net/mail"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"selfhosted/html"
	"selfhosted/toast"
	"time"

	"github.com/go-chi/httprate"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var loginRateLimiter = httprate.NewRateLimiter(5, time.Minute)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	sess, ok := r.Context().Value(app.SessionKey).(store.GetSessionByUuidRow)
	if ok && sess.ID > 0 {
		slog.Error("User already logged in", "session_id", sess.ID, "user_id", sess.UserID)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	html.LoginPage().Render(r.Context(), w)
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("email") == "" || r.FormValue("password") == "" {
		w.WriteHeader(http.StatusBadRequest)
		html.LoginForm().Render(r.Context(), w)
		return
	}

	email := r.FormValue("email")

	if loginRateLimiter.OnLimit(w, r, email) {
		w.WriteHeader(http.StatusTooManyRequests)
		html.LoginForm().Render(r.Context(), w)
		toast.Error("Too many requests", "You have exceeded the maximum number of login attempts. Please try again later.").Send(r.Context(), w)
		return
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		slog.Error("Invalid email format", "email", email, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		html.LoginForm().Render(r.Context(), w)
		toast.Error("Login failed", "The credentials you provided are invalid.").Send(r.Context(), w)
		return
	}

	u, err := store.New(database.DB).GetUserByEmail(r.Context(), email)
	if err != nil || u.ID == 0 {
		slog.Error("User not found", "email", email, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		html.LoginForm().Render(r.Context(), w)
		toast.Error("Login failed", "The credentials you provided are invalid.").Send(r.Context(), w)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(r.FormValue("password")))
	if err != nil {
		slog.Error("Password mismatch", "email", email, "error", err)
		w.WriteHeader(http.StatusBadRequest)
		html.LoginForm().Render(r.Context(), w)
		toast.Error("Login failed", "The credentials you provided are invalid.").Send(r.Context(), w)
		return
	}

	uuid := uuid.NewString()
	expiresAt := time.Now().Add(24 * time.Hour * 30).UTC()

	err = store.New(database.DB).CreateSession(r.Context(), store.CreateSessionParams{
		Uuid:      uuid,
		UserID:    u.ID,
		ExpiresAt: expiresAt,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		HttpOnly: true,
		Value:    uuid,
		Expires:  expiresAt,
		SameSite: http.SameSiteLaxMode,
	})

	loginRateLimiter.Counter().IncrementBy(email, time.Now(), -1)

	// w.Header().Set("HX-Location", "/")
	w.Header().Set("HX-Trigger", `{"redirect": "/"}`)
	w.Header().Set("HX-Push-Url", `/`)
	html.LoginForm().Render(r.Context(), w)
	toast.Success("Welcome back!", "You have successfully logged in.").Send(r.Context(), w)
}

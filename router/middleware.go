package router

import (
	"context"
	"log/slog"
	"net/http"
	"selfhosted/app"
	"selfhosted/database"
	"selfhosted/database/store"
	"time"
)

func SetupMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		if app.AdminCount == 0 {
			http.Redirect(w, r, "/setup", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func UserMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		if err != nil || c.Value == "" {
			slog.Error("Failed to get session cookie", "error", err)
			next.ServeHTTP(w, r)
			return
		}

		s, err := store.New(database.DB).GetSessionByUuid(r.Context(), c.Value)
		if err != nil || s.ID == 0 {
			slog.Error("Failed to get session by UUID", "error", err, "cookie_value", c.Value)
			next.ServeHTTP(w, r)
			return
		}

		if s.ExpiresAt.Before(time.Now()) {
			slog.Debug("Session expired", "session_id", s.ID, "user_id", s.UserID)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), app.SessionKey, s)
		slog.Debug("Context updated with session", "session_id", s.ID, "user_id", s.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess, ok := r.Context().Value(app.SessionKey).(store.GetSessionByUuidRow)
		if !ok || sess.ID == 0 {
			slog.Error("No valid session found in context", "context_value", r.Context().Value(app.SessionKey))
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func SetCurrentUrlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), app.CurrentUrlKey, r.URL.Path)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

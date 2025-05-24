package router

import (
	"net/http"
	"template/app"
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

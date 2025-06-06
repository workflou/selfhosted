package router

import (
	"net/http"
	"selfhosted/handler"
	"selfhosted/html"
	"selfhosted/static"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(time.Second * 30))
	r.Use(middleware.GetHead)

	r.Group(func(r chi.Router) {
		r.Use(middleware.NoCache)

		r.Group(func(r chi.Router) {
			r.Use(UserMiddleware)
			r.Use(SetCurrentUrlMiddleware)

			r.Group(func(r chi.Router) {
				r.Use(SetupMiddleware)

				r.Post("/setup", handler.SetupForm)
				r.Get("/login", handler.LoginPage)
				r.Post("/login", handler.LoginForm)
			})

			r.Group(func(r chi.Router) {
				r.Use(SetupMiddleware)
				r.Use(AuthMiddleware)

				r.Get("/", handler.HomePage)
				r.Get("/about", templ.Handler(html.AboutPage()).ServeHTTP)
				r.Get("/logout", handler.Logout)
				r.Get("/settings", handler.SettingsPage)
				r.Post("/settings/name", handler.SettingsNameForm)
				r.Post("/settings/avatar", handler.SettingsAvatarForm)
			})
		})

		r.Get("/setup", handler.SetupPage)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.SetHeader("Cache-Control", "public, max-age=31536000"))

		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))
		r.Handle("/uploads/avatars/*", http.StripPrefix("/uploads/avatars/", http.FileServer(http.Dir("./uploads/avatars"))))
	})

	return r
}

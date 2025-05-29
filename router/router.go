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
	r.Use(middleware.NoCache)

	r.Group(func(r chi.Router) {
		r.Use(UserMiddleware)

		r.Group(func(r chi.Router) {
			r.Use(SetupMiddleware)

			r.Post("/setup", handler.SetupForm)
			r.Get("/login", handler.LoginPage)
			r.Post("/login", handler.LoginForm)
		})

		r.Group(func(r chi.Router) {
			r.Use(SetupMiddleware)
			r.Use(AuthMiddleware)

			r.Get("/", templ.Handler(html.HomePage()).ServeHTTP)
		})
	})

	r.Get("/setup", handler.SetupPage)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return r
}

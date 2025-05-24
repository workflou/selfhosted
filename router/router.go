package router

import (
	"net/http"
	"template/html"
	"template/static"
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

	r.Use(SetupMiddleware)

	r.Get("/", templ.Handler(html.HomePage()).ServeHTTP)

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(static.FS))))

	return r
}

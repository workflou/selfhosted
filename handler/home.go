package handler

import (
	"net/http"
	"selfhosted/html"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	html.HomePage().Render(r.Context(), w)
}

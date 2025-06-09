package handler

import (
	"net/http"
	"selfhosted/html"
)

func UserModal(w http.ResponseWriter, r *http.Request) {
	html.UserModal().Render(r.Context(), w)
}

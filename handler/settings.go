package handler

import (
	"net/http"
	"selfhosted/html"
)

func SettingsPage(w http.ResponseWriter, r *http.Request) {
	html.SettingsPage().Render(r.Context(), w)
}

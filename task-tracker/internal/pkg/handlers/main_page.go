package handlers

import "net/http"

func (h *Service) MainPage(w http.ResponseWriter, r *http.Request) {
	outputHTML(w, r, "static/main.html")
}

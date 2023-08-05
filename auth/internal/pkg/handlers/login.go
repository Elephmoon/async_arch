package handlers

import (
	"auth/internal/pkg/models"
	"encoding/json"
	"net/http"
)

func (h *Service) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		outputHTML(w, r, "static/login.html")
		return
	}

	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	userPass := r.Form.Get("password")
	userName := r.Form.Get("username")

	token, err := h.userUseCase.LoginUser(r.Context(), userName, userPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("content-type", "application/json")
	authResp := models.AuthResponse{Token: token}
	bodyResp, err := json.Marshal(authResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(bodyResp)
}

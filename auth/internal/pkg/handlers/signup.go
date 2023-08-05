package handlers

import (
	"auth/internal/pkg/models"
	"encoding/json"
	"net/http"
)

func (h *Service) Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		outputHTML(w, r, "static/signup.html")
		return
	}

	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	userName := r.Form.Get("username")
	userPassword := r.Form.Get("password")
	userRole := r.Form.Get("role")

	token, err := h.userUseCase.CreateUser(r.Context(), userName, userRole, userPassword)
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

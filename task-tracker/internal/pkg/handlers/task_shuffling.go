package handlers

import (
	"net/http"
	"strings"
	"task-tracker/internal/pkg/jwt"
)

func (h *Service) Shuffle(w http.ResponseWriter, r *http.Request) {
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	jwtToken := r.Form.Get("jwt-token")
	if strings.TrimSpace(jwtToken) == "" {
		http.Error(w, "token is empty", http.StatusBadRequest)
		return
	}
	validateRes, err := jwt.ValidateToken(jwtToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if !validateRes.Valid {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	err = h.tasksUseCase.Shuffle(r.Context(), validateRes.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

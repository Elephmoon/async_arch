package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"task-tracker/internal/pkg/jwt"
)

func (h *Service) GetMyTasks(w http.ResponseWriter, r *http.Request) {
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

	tasks, err := h.tasksUseCase.GetUserTasks(r.Context(), validateRes.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	jsonBody, err := json.Marshal(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBody)
}

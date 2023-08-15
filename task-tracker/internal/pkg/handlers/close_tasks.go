package handlers

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"task-tracker/internal/pkg/jwt"
	"task-tracker/internal/pkg/models"
)

func (h *Service) CloseTask(w http.ResponseWriter, r *http.Request) {
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

	taskID := r.Form.Get("task-id")
	if strings.TrimSpace(jwtToken) == "" {
		http.Error(w, "task_id is empty", http.StatusBadRequest)
		return
	}
	taskGuid, err := uuid.Parse(taskID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.tasksUseCase.CloseTask(r.Context(), taskGuid, validateRes.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	resp := models.TaskResponse{TaskID: taskGuid}
	jsonBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(jsonBody)
}

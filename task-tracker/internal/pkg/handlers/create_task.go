package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"task-tracker/internal/pkg/jwt"
	"task-tracker/internal/pkg/models"
)

func (h *Service) CreateTask(w http.ResponseWriter, r *http.Request) {
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

	taskName := r.Form.Get("task-name")
	if strings.TrimSpace(taskName) == "" {
		http.Error(w, "taskName is empty", http.StatusBadRequest)
		return
	}
	taskDescription := r.Form.Get("task-description")
	if strings.TrimSpace(taskDescription) == "" {
		http.Error(w, "taskDescription is empty", http.StatusBadRequest)
		return
	}
	assigneeUserName := r.Form.Get("assignee-username")
	if strings.TrimSpace(assigneeUserName) == "" {
		http.Error(w, "taskDescription is empty", http.StatusBadRequest)
		return
	}

	taskID, err := h.tasksUseCase.CreateTask(r.Context(), assigneeUserName, taskName, taskDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Add("content-type", "application/json")
	resp := models.TaskResponse{TaskID: taskID}
	jsonBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write(jsonBody)
}

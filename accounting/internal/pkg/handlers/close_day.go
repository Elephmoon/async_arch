package handlers

import "net/http"

func (h *Service) CloseDay(w http.ResponseWriter, r *http.Request) {
	err := h.dayCloser.CloseDay(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

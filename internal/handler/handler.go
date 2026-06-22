package handler

import (
	"net/http"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"status": "ok"}
	WriteJSON(w, http.StatusOK, data)
}

package handlers

import (
	"awesomeProject1/internal/store"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	store *store.Store
}

func New(s *store.Store) *Handler {
	return &Handler{store: s}
}

type errResp struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func parseID(r *http.Request) (int, bool) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		return 0, false
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, true // "true" = id был, но он невалидный
	}
	return id, false
}

func (h *Handler) Tasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTasks(w, r)
	case http.MethodPost:
		h.createTask(w, r)
	case http.MethodPatch:
		h.patchTask(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, errResp{Error: "method not allowed"})
	}
}

func (h *Handler) getTasks(w http.ResponseWriter, r *http.Request) {
	id, invalid := parseID(r)
	if invalid {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid id"})
		return
	}

	// GET /tasks?id=...
	if id != 0 {
		t, ok := h.store.Get(id)
		if !ok {
			writeJSON(w, http.StatusNotFound, errResp{Error: "task not found"})
			return
		}
		writeJSON(w, http.StatusOK, t)
		return
	}

	// GET /tasks
	all := h.store.List()
	writeJSON(w, http.StatusOK, all)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Title string `json:"title"`
	}
	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid title"})
		return
	}

	title := strings.TrimSpace(body.Title)
	if title == "" {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid title"})
		return
	}

	t := h.store.Create(title)
	writeJSON(w, http.StatusCreated, t)
}

func (h *Handler) patchTask(w http.ResponseWriter, r *http.Request) {
	id, invalid := parseID(r)
	if invalid || id == 0 {
		// в PATCH id обязателен
		writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid id"})
		return
	}

	type req struct {
		Done *bool `json:"done"`
	}
	var body req
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Done == nil {
		writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid done"})
		return
	}

	ok := h.store.UpdateDone(id, *body.Done)
	if !ok {
		writeJSON(w, http.StatusNotFound, errResp{Error: "task not found"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"updated": true})
}

package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"golang/internal/usecase"
	"golang/pkg/modules"

	"github.com/go-chi/chi/v5"
)

type UsersHandler struct {
	uc usecase.UserService
}

func NewUsersHandler(uc usecase.UserService) *UsersHandler {
	return &UsersHandler{uc: uc}
}

func (h *UsersHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.uc.GetUsers(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *UsersHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	user, err := h.uc.GetUserByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, user)
}

type createUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func (h *UsersHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	created, err := h.uc.CreateUser(r.Context(), modules.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	var patch usecase.UpdateUserPatch
	if err := decodeJSON(r, &patch); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	updated, err := h.uc.UpdateUser(r.Context(), id, patch)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}

	rows, err := h.uc.DeleteUser(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"deleted":      true,
		"rowsAffected": rows,
	})
}

/* helpers */

func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	raw := chi.URLParam(r, "id")
	id, err := strconv.Atoi(raw)
	if err != nil || id <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	msg := "internal server error"

	switch {
	case errors.Is(err, modules.ErrInvalidInput):
		status = http.StatusBadRequest
		msg = err.Error()
	case errors.Is(err, modules.ErrUserNotFound):
		status = http.StatusNotFound
		msg = err.Error()
	case errors.Is(err, modules.ErrConflict):
		status = http.StatusConflict
		msg = err.Error()
	}

	if status == http.StatusInternalServerError {
		log.Printf("internal error: %v", err)
	}

	writeJSON(w, status, map[string]any{"error": msg})
}

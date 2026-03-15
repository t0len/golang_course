package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"practice5/internal/models"
	"practice5/internal/repository"

	"github.com/google/uuid"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /users?page=1&page_size=10&order_by=name&name=alice&gender=male
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	pageSize, _ := strconv.Atoi(q.Get("page_size"))

	params := models.FilterParams{
		ID:        q.Get("id"),
		Name:      q.Get("name"),
		Email:     q.Get("email"),
		Gender:    q.Get("gender"),
		BirthDate: q.Get("birth_date"),
		OrderBy:   q.Get("order_by"),
		Page:      page,
		PageSize:  pageSize,
	}

	result, err := h.repo.GetPaginatedUsers(params)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

// GET /users/common-friends?user1=<uuid>&user2=<uuid>
func (h *Handler) GetCommonFriends(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	id1, err := uuid.Parse(q.Get("user1"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user1 uuid")
		return
	}
	id2, err := uuid.Parse(q.Get("user2"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid user2 uuid")
		return
	}
	if id1 == id2 {
		writeError(w, http.StatusBadRequest, "user1 and user2 must be different")
		return
	}

	friends, err := h.repo.GetCommonFriends(id1, id2)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"common_friends": friends, "count": len(friends)})
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users", h.GetUsers)
	mux.HandleFunc("/users/common-friends", h.GetCommonFriends)
}

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jtorre/qisurChallenge/internal/models"
	"github.com/jtorre/qisurChallenge/internal/repository"
	"github.com/jtorre/qisurChallenge/internal/ws"
)

type CategoryHandler struct {
	repo *repository.CategoryRepository
	hub  *ws.Hub
}

func NewCategoryHandler(repo *repository.CategoryRepository, hub *ws.Hub) *CategoryHandler {
	return &CategoryHandler{repo: repo, hub: hub}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *CategoryHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	category, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if cat.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	created, err := h.repo.Create(r.Context(), &cat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.hub.Broadcast("category_created", created)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	var cat models.Category
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if cat.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	updated, err := h.repo.Update(r.Context(), id, &cat)
	if err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	h.hub.Broadcast("category_updated", updated)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	err = h.repo.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	h.hub.Broadcast("category_deleted", map[string]string{"id": id.String()})

	w.WriteHeader(http.StatusNoContent)
}

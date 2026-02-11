package rest

import (
	"encoding/json"
	"net/http"

	"github.com/daniil-dev/project-store/backend/services/catalog/internal/domain"
	"github.com/daniil-dev/project-store/backend/services/catalog/internal/repository"
)

type Handler struct {
	repo *repository.ProductRepository
}

func NewHandler(repo *repository.ProductRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var p domain.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "Bad JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.CreateProduct(r.Context(), &p); err != nil {
		http.Error(w, "DB Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

func (h *Handler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	products, err := h.repo.GetAllProducts(r.Context())
	if err != nil {
		http.Error(w, "DB Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

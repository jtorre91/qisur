package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jtorre/qisurChallenge/internal/repository"
	"github.com/jtorre/qisurChallenge/internal/utils"
)

type SearchHandler struct {
	repo *repository.SearchRepository
}

func NewSearchHandler(repo *repository.SearchRepository) *SearchHandler {
	return &SearchHandler{repo: repo}
}

func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	searchType := r.URL.Query().Get("type")

	if searchType == "" {
		http.Error(w, "type parameter is required (product or category)", http.StatusBadRequest)
		return
	}

	searchType = strings.ToLower(searchType)

	if searchType == "product" {
		h.searchProducts(w, r)
	} else if searchType == "category" {
		h.searchCategories(w, r)
	} else {
		http.Error(w, "invalid type. Use 'product' or 'category'", http.StatusBadRequest)
	}
}

func (h *SearchHandler) searchProducts(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	sortBy := utils.ExtractSortByParam(r)
	order := utils.ExtractOrderParam(r)
	page := utils.ExtractPageParam(r)
	limit := utils.ExtractLimitParam(r)
	minPrice := utils.ExtractMinPriceParam(r)
	maxPrice := utils.ExtractMaxPriceParam(r)

	if minPrice != nil && maxPrice != nil && *minPrice > *maxPrice {
		http.Error(w, "min_price cannot be greater than max_price", http.StatusBadRequest)
		return
	}

	params := repository.SearchProductsParams{
		Query:    query,
		MinPrice: minPrice,
		MaxPrice: maxPrice,
		SortBy:   sortBy,
		Order:    order,
		Page:     page,
		Limit:    limit,
	}

	result, err := h.repo.SearchProducts(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (h *SearchHandler) searchCategories(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	page := utils.ExtractPageParam(r)
	limit := utils.ExtractLimitParam(r)
	sortBy := utils.ExtractSortByParam(r)
	order := utils.ExtractOrderParam(r)

	params := repository.SearchCategoriesParams{
		Query:  query,
		SortBy: sortBy,
		Order:  order,
		Page:   page,
		Limit:  limit,
	}

	result, err := h.repo.SearchCategories(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

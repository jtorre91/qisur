package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jtorre/qisurChallenge/internal/models"
)

type SearchRepository struct {
	pool *pgxpool.Pool
}

type SearchProductsParams struct {
	Query    string
	MinPrice *float64
	MaxPrice *float64
	SortBy   string // name, price, stock, created_at
	Order    string // ASC, DESC
	Page     int
	Limit    int
}

type SearchCategoriesParams struct {
	Query  string
	SortBy string // name, created_at
	Order  string // ASC, DESC
	Page   int
	Limit  int
}

type SearchResult struct {
	Items      interface{} `json:"items"`
	Total      int         `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"total_pages"`
}

func NewSearchRepository(pool *pgxpool.Pool) *SearchRepository {
	return &SearchRepository{pool: pool}
}

func (r *SearchRepository) SearchProducts(ctx context.Context, params SearchProductsParams) (*SearchResult, error) {
	// Validate and set defaults
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Page < 1 {
		params.Page = 1
	}

	// Validate sort_by
	validSortByFields := map[string]bool{
		"name":       true,
		"price":      true,
		"stock":      true,
		"created_at": true,
	}
	if _, valid := validSortByFields[params.SortBy]; !valid {
		params.SortBy = "created_at"
	}

	// Validate order
	if strings.ToUpper(params.Order) != "ASC" && strings.ToUpper(params.Order) != "DESC" {
		params.Order = "DESC"
	}
	params.Order = strings.ToUpper(params.Order)

	offset := (params.Page - 1) * params.Limit

	// Build query
	query := "SELECT id, name, description, price, stock, created_at, updated_at FROM products WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Search by query
	if params.Query != "" {
		searchPattern := fmt.Sprintf("%%%s%%", params.Query)
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Price filters
	if params.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argIndex)
		args = append(args, *params.MinPrice)
		argIndex++
	}
	if params.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argIndex)
		args = append(args, *params.MaxPrice)
		argIndex++
	}

	// Count total before pagination
	countQuery := strings.Replace(query, "SELECT id, name, description, price, stock, created_at, updated_at", "SELECT COUNT(*)", 1)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Add sorting and pagination
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", params.SortBy, params.Order, argIndex, argIndex+1)
	args = append(args, params.Limit, offset)

	// Execute query
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search products: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var prod models.Product
		err := rows.Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, prod)
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &SearchResult{
		Items:      products,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

func (r *SearchRepository) SearchCategories(ctx context.Context, params SearchCategoriesParams) (*SearchResult, error) {
	// Validate and set defaults
	if params.Limit == 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}
	if params.Page < 1 {
		params.Page = 1
	}

	// Validate sort_by
	validSortByFields := map[string]bool{
		"name":       true,
		"created_at": true,
	}
	if _, valid := validSortByFields[params.SortBy]; !valid {
		params.SortBy = "created_at"
	}

	// Validate order
	if strings.ToUpper(params.Order) != "ASC" && strings.ToUpper(params.Order) != "DESC" {
		params.Order = "DESC"
	}
	params.Order = strings.ToUpper(params.Order)

	offset := (params.Page - 1) * params.Limit

	// Build query
	query := "SELECT id, name, description, created_at, updated_at FROM categories WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Search by query
	if params.Query != "" {
		searchPattern := fmt.Sprintf("%%%s%%", params.Query)
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d)", argIndex, argIndex+1)
		args = append(args, searchPattern, searchPattern)
		argIndex += 2
	}

	// Count total before pagination
	countQuery := strings.Replace(query, "SELECT id, name, description, created_at, updated_at", "SELECT COUNT(*)", 1)
	var total int
	err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to count categories: %w", err)
	}

	// Add sorting and pagination
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", params.SortBy, params.Order, argIndex, argIndex+1)
	args = append(args, params.Limit, offset)

	// Execute query
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		err := rows.Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, cat)
	}

	totalPages := (total + params.Limit - 1) / params.Limit

	return &SearchResult{
		Items:      categories,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

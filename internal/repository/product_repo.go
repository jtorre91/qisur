package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jtorre/qisurChallenge/internal/models"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
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

	return products, rows.Err()
}

func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var prod models.Product
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`, id).Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &prod, nil
}

func (r *ProductRepository) Create(ctx context.Context, prod *models.Product) (*models.Product, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	prod.ID = uuid.New()

	err = tx.QueryRow(ctx, `
		INSERT INTO products (id, name, description, price, stock)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, price, stock, created_at, updated_at
	`, prod.ID, prod.Name, prod.Description, prod.Price, prod.Stock).Scan(
		&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Insert into product_history for initial record
	_, err = tx.Exec(ctx, `
		INSERT INTO product_history (product_id, price, stock)
		VALUES ($1, $2, $3)
	`, prod.ID, prod.Price, prod.Stock)
	if err != nil {
		return nil, fmt.Errorf("failed to create product history: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return prod, nil
}

func (r *ProductRepository) Update(ctx context.Context, id uuid.UUID, prod *models.Product) (*models.Product, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get old product to check if price or stock changed
	var oldPrice float64
	var oldStock int
	err = tx.QueryRow(ctx, `
		SELECT price, stock FROM products WHERE id = $1
	`, id).Scan(&oldPrice, &oldStock)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Update product
	err = tx.QueryRow(ctx, `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, updated_at = now()
		WHERE id = $5
		RETURNING id, name, description, price, stock, created_at, updated_at
	`, prod.Name, prod.Description, prod.Price, prod.Stock, id).Scan(
		&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Stock, &prod.CreatedAt, &prod.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Record history if price or stock changed
	if prod.Price != oldPrice || prod.Stock != oldStock {
		_, err = tx.Exec(ctx, `
			INSERT INTO product_history (product_id, price, stock)
			VALUES ($1, $2, $3)
		`, id, prod.Price, prod.Stock)
		if err != nil {
			return nil, fmt.Errorf("failed to create product history: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return prod, nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM products WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("product not found")
	}

	return nil
}

func (r *ProductRepository) GetHistory(ctx context.Context, productID uuid.UUID, startDate *string, endDate *string, limit int, offset int) ([]models.ProductHistory, error) {
	query := `
		SELECT id, product_id, price, stock, changed_at
		FROM product_history
		WHERE product_id = $1
	`
	args := []interface{}{productID}
	argIndex := 2

	// Add date range filters if provided
	if startDate != nil && *startDate != "" {
		query += fmt.Sprintf(" AND changed_at::date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil && *endDate != "" {
		query += fmt.Sprintf(" AND changed_at::date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY changed_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query product history: %w", err)
	}
	defer rows.Close()

	history := make([]models.ProductHistory, 0)
	for rows.Next() {
		var h models.ProductHistory
		err := rows.Scan(&h.ID, &h.ProductID, &h.Price, &h.Stock, &h.ChangedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history: %w", err)
		}
		history = append(history, h)
	}

	return history, rows.Err()
}

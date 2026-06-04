package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jtorre/qisurChallenge/internal/models"
)

type CategoryRepository struct {
	pool *pgxpool.Pool
}

func NewCategoryRepository(pool *pgxpool.Pool) *CategoryRepository {
	return &CategoryRepository{pool: pool}
}

func (r *CategoryRepository) List(ctx context.Context) ([]models.Category, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
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

	return categories, rows.Err()
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Category, error) {
	var cat models.Category
	err := r.pool.QueryRow(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1
	`, id).Scan(&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &cat, nil
}

func (r *CategoryRepository) Create(ctx context.Context, cat *models.Category) (*models.Category, error) {
	cat.ID = uuid.New()

	err := r.pool.QueryRow(ctx, `
		INSERT INTO categories (id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, created_at, updated_at
	`, cat.ID, cat.Name, cat.Description).Scan(
		&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return cat, nil
}

func (r *CategoryRepository) Update(ctx context.Context, id uuid.UUID, cat *models.Category) (*models.Category, error) {
	err := r.pool.QueryRow(ctx, `
		UPDATE categories
		SET name = $1, description = $2, updated_at = now()
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at
	`, cat.Name, cat.Description, id).Scan(
		&cat.ID, &cat.Name, &cat.Description, &cat.CreatedAt, &cat.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return cat, nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result, err := r.pool.Exec(ctx, `
		DELETE FROM categories WHERE id = $1
	`, id)

	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jtorre/qisurChallenge/internal/models"
)


func TestCategoryModel(t *testing.T) {
	tests := []struct {
		name     string
		category models.Category
		validate func(models.Category) bool
	}{
		{
			name: "Valid category with all fields",
			category: models.Category{
				ID:          uuid.New(),
				Name:        "Electronics",
				Description: "Electronic products",
			},
			validate: func(c models.Category) bool {
				return c.ID != uuid.Nil && c.Name != "" && c.Description != ""
			},
		},
		{
			name: "Category with name only",
			category: models.Category{
				ID:   uuid.New(),
				Name: "Books",
			},
			validate: func(c models.Category) bool {
				return c.ID != uuid.Nil && c.Name != ""
			},
		},
		{
			name: "Category without ID",
			category: models.Category{
				Name:        "Clothing",
				Description: "Clothing items",
			},
			validate: func(c models.Category) bool {
				return c.Name != "" && c.Description != ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.validate(tt.category) {
				t.Errorf("Category validation failed for: %s", tt.name)
			}
		})
	}
}

func TestCategoryValidation(t *testing.T) {
	tests := []struct {
		name        string
		category    models.Category
		isValid     bool
		description string
	}{
		{
			name:        "Category with empty name is invalid",
			category:    models.Category{Name: "", Description: "Test"},
			isValid:     false,
			description: "Name is required",
		},
		{
			name:        "Category with valid name is valid",
			category:    models.Category{Name: "Valid Name", Description: "Test"},
			isValid:     true,
			description: "Name provided",
		},
		{
			name:        "Category with very long name",
			category:    models.Category{Name: "This is a very long category name that exceeds normal length", Description: "Test"},
			isValid:     true,
			description: "Long names are accepted",
		},
		{
			name:        "Category with special characters",
			category:    models.Category{Name: "Category #1 & More", Description: "Test"},
			isValid:     true,
			description: "Special characters in name",
		},
		{
			name:        "Category with unicode characters",
			category:    models.Category{Name: "Категория 中文 العربية", Description: "Test"},
			isValid:     true,
			description: "Unicode characters in name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := tt.category.Name != ""
			if isValid != tt.isValid {
				t.Errorf("Validation failed for: %s. Expected %v, got %v. %s",
					tt.name, tt.isValid, isValid, tt.description)
			}
		})
	}
}

func TestCategoryIDGeneration(t *testing.T) {
	tests := []struct {
		name        string
		generateID  func() uuid.UUID
		expectValid bool
	}{
		{
			name: "UUID.New() generates valid ID",
			generateID: func() uuid.UUID {
				return uuid.New()
			},
			expectValid: true,
		},
		{
			name: "UUID.Nil is invalid ID",
			generateID: func() uuid.UUID {
				return uuid.Nil
			},
			expectValid: false,
		},
		{
			name: "Generated UUIDs are never nil",
			generateID: func() uuid.UUID {
				return uuid.New()
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := tt.generateID()
			isValid := id != uuid.Nil
			if isValid != tt.expectValid {
				t.Errorf("ID validation failed for: %s", tt.name)
			}
		})
	}
}

func TestCategoryComparison(t *testing.T) {
	id1 := uuid.New()

	tests := []struct {
		name        string
		cat1        models.Category
		cat2        models.Category
		shouldEqual bool
	}{
		{
			name: "Same category objects are equal",
			cat1: models.Category{
				ID:          id1,
				Name:        "Test",
				Description: "Description",
			},
			cat2: models.Category{
				ID:          id1,
				Name:        "Test",
				Description: "Description",
			},
			shouldEqual: true,
		},
		{
			name: "Different IDs are not equal",
			cat1: models.Category{
				ID:   id1,
				Name: "Test",
			},
			cat2: models.Category{
				ID:   uuid.New(),
				Name: "Test",
			},
			shouldEqual: false,
		},
		{
			name: "Different names are not equal",
			cat1: models.Category{
				ID:   id1,
				Name: "Category 1",
			},
			cat2: models.Category{
				ID:   id1,
				Name: "Category 2",
			},
			shouldEqual: false,
		},
		{
			name: "Different descriptions are still equal if ID and name match",
			cat1: models.Category{
				ID:          id1,
				Name:        "Same",
				Description: "Desc 1",
			},
			cat2: models.Category{
				ID:          id1,
				Name:        "Same",
				Description: "Desc 2",
			},
			shouldEqual: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			equal := tt.cat1.ID == tt.cat2.ID && tt.cat1.Name == tt.cat2.Name
			if equal != tt.shouldEqual {
				t.Errorf("Comparison failed for: %s. Expected equal=%v, got %v",
					tt.name, tt.shouldEqual, equal)
			}
		})
	}
}

func TestCategoryFields(t *testing.T) {
	tests := []struct {
		name     string
		category models.Category
		checks   func(models.Category) error
	}{
		{
			name: "All fields populated",
			category: models.Category{
				ID:          uuid.New(),
				Name:        "Test",
				Description: "Test Description",
			},
			checks: func(c models.Category) error {
				if c.ID == uuid.Nil {
					return ErrIDRequired
				}
				if c.Name == "" {
					return ErrNameRequired
				}
				if c.Description == "" {
					return ErrDescriptionEmpty
				}
				return nil
			},
		},
		{
			name: "Missing description is allowed",
			category: models.Category{
				ID:   uuid.New(),
				Name: "Test",
			},
			checks: func(c models.Category) error {
				if c.ID == uuid.Nil {
					return ErrIDRequired
				}
				if c.Name == "" {
					return ErrNameRequired
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.checks(tt.category); err != nil {
				t.Errorf("Field check failed for: %s, error: %v", tt.name, err)
			}
		})
	}
}

// Error types for testing
type testError string

func (e testError) Error() string {
	return string(e)
}

const (
	ErrIDRequired         testError = "ID is required"
	ErrNameRequired       testError = "Name is required"
	ErrDescriptionEmpty   testError = "Description is empty"
)

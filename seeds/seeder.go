package seeds

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func Run(pool *pgxpool.Pool) error {
	ctx := context.Background()

	// Create users
	adminID := uuid.New()
	clientID := uuid.New()

	adminHash, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	clientHash, _ := bcrypt.GenerateFromPassword([]byte("client123"), bcrypt.DefaultCost)

	_, err := pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO NOTHING
	`, adminID, "admin@qisur.com", string(adminHash), "admin")
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	_, err = pool.Exec(ctx, `
		INSERT INTO users (id, email, password_hash, role)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (email) DO NOTHING
	`, clientID, "client@qisur.com", string(clientHash), "client")
	if err != nil {
		return fmt.Errorf("failed to insert client user: %w", err)
	}

	// Create categories
	categories := []struct {
		id   uuid.UUID
		name string
		desc string
	}{
		{uuid.New(), "Lácteos", "Leches, yogures, quesos y derivados lácteos"},
		{uuid.New(), "Carnes y Aves", "Carnes frescas y aves de primera calidad"},
		{uuid.New(), "Frutas y Verduras", "Productos frescos de temporada"},
		{uuid.New(), "Bebidas", "Bebidas variadas: refrescos, jugos, agua"},
		{uuid.New(), "Panadería", "Pan, facturas y productos de panadería"},
		{uuid.New(), "Almacén", "Productos de almacén: arroz, pasta, aceite, azúcar"},
		{uuid.New(), "Congelados", "Productos congelados listos para cocinar"},
		{uuid.New(), "Limpieza", "Productos de limpieza y cuidado del hogar"},
	}

	for _, cat := range categories {
		_, err := pool.Exec(ctx, `
			INSERT INTO categories (id, name, description)
			VALUES ($1, $2, $3)
			ON CONFLICT DO NOTHING
		`, cat.id, cat.name, cat.desc)
		if err != nil {
			return fmt.Errorf("failed to insert category: %w", err)
		}
	}

	// Create products (Supermercado)
	products := []struct {
		name  string
		desc  string
		price float64
		stock int
		cats  []uuid.UUID
	}{
		// Lácteos
		{
			"Leche Entera 1L", "Leche fresca entera de la mejor calidad", 1.50, 150,
			[]uuid.UUID{categories[0].id},
		},
		{
			"Yogur Natural 500g", "Yogur natural sin azúcar agregada", 2.99, 80,
			[]uuid.UUID{categories[0].id},
		},
		{
			"Queso Fresco 400g", "Queso fresco artesanal de la región", 5.50, 40,
			[]uuid.UUID{categories[0].id},
		},
		{
			"Mantequilla 200g", "Mantequilla pura sin conservantes", 4.25, 60,
			[]uuid.UUID{categories[0].id},
		},
		// Carnes y Aves
		{
			"Pechuga de Pollo 1kg", "Pechuga de pollo fresca y magra", 7.99, 50,
			[]uuid.UUID{categories[1].id},
		},
		{
			"Carne Picada 500g", "Carne molida premium para todas tus recetas", 6.50, 45,
			[]uuid.UUID{categories[1].id},
		},
		{
			"Costilla de Cerdo 1kg", "Costillas frescas de cerdo de primera", 8.99, 30,
			[]uuid.UUID{categories[1].id},
		},
		// Frutas y Verduras
		{
			"Manzana Roja x6", "Manzanas rojas frescas por peso", 3.99, 100,
			[]uuid.UUID{categories[2].id},
		},
		{
			"Plátano x6", "Plátanos maduros listos para comer", 2.49, 120,
			[]uuid.UUID{categories[2].id},
		},
		{
			"Lechuga Mantecosa", "Lechuga fresca y crujiente para ensaladas", 1.99, 90,
			[]uuid.UUID{categories[2].id},
		},
		{
			"Tomate Rojo 1kg", "Tomates frescos recién cosechados", 2.75, 110,
			[]uuid.UUID{categories[2].id},
		},
		// Bebidas
		{
			"Gaseosa 2L", "Gaseosa refrescante de colores variados", 1.99, 200,
			[]uuid.UUID{categories[3].id},
		},
		{
			"Jugo Natural 1L", "Jugo 100% natural sin azúcar agregada", 3.50, 70,
			[]uuid.UUID{categories[3].id},
		},
		{
			"Agua Mineral 1.5L", "Agua mineral pura y cristalina", 0.99, 300,
			[]uuid.UUID{categories[3].id},
		},
		{
			"Cerveza Pack 6", "Cerveza artesanal premium", 8.99, 60,
			[]uuid.UUID{categories[3].id},
		},
		// Panadería
		{
			"Pan de Molde", "Pan de molde integral recién horneado", 2.50, 80,
			[]uuid.UUID{categories[4].id},
		},
		{
			"Facturas Surtidas", "Caja de 6 facturas variadas", 4.99, 50,
			[]uuid.UUID{categories[4].id},
		},
		{
			"Biscochos de Grasa", "Biscochos de grasa tradicionales", 1.75, 100,
			[]uuid.UUID{categories[4].id},
		},
		// Almacén
		{
			"Aceite de Oliva 500ml", "Aceite de oliva virgen extra premium", 6.99, 40,
			[]uuid.UUID{categories[5].id},
		},
		{
			"Arroz Blanco 1kg", "Arroz de grano largo de la mejor cosecha", 2.25, 120,
			[]uuid.UUID{categories[5].id},
		},
		{
			"Pasta Seca 500g", "Pasta variada: fideos, moños, penne", 1.50, 150,
			[]uuid.UUID{categories[5].id},
		},
		{
			"Azúcar 1kg", "Azúcar blanca refinada", 1.99, 100,
			[]uuid.UUID{categories[5].id},
		},
		{
			"Sal Fina 1kg", "Sal refinada yodada para cocina", 0.75, 200,
			[]uuid.UUID{categories[5].id},
		},
		// Congelados
		{
			"Pizza Congelada 400g", "Pizza congelada lista para hornear", 4.99, 60,
			[]uuid.UUID{categories[6].id},
		},
		{
			"Helado 1L", "Helado variado de sabores", 5.50, 45,
			[]uuid.UUID{categories[6].id},
		},
		{
			"Verduras Congeladas", "Mix de verduras congeladas para cocinar", 3.25, 70,
			[]uuid.UUID{categories[6].id},
		},
		// Limpieza
		{
			"Detergente Líquido 1L", "Detergente concentrado para ropa", 3.75, 80,
			[]uuid.UUID{categories[7].id},
		},
		{
			"Jabón de Pisos 500ml", "Jabón desengrasante para limpieza", 2.50, 100,
			[]uuid.UUID{categories[7].id},
		},
		{
			"Desinfectante Spray", "Desinfectante multiusos de 500ml", 4.25, 90,
			[]uuid.UUID{categories[7].id},
		},
	}

	for _, prod := range products {
		prodID := uuid.New()
		_, err := pool.Exec(ctx, `
			INSERT INTO products (id, name, description, price, stock)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT DO NOTHING
		`, prodID, prod.name, prod.desc, prod.price, prod.stock)
		if err != nil {
			return fmt.Errorf("failed to insert product: %w", err)
		}

		// Insert product-category
		for _, catID := range prod.cats {
			_, err := pool.Exec(ctx, `
				INSERT INTO product_category (product_id, category_id)
				VALUES ($1, $2)
				ON CONFLICT DO NOTHING
			`, prodID, catID)
			if err != nil {
				return fmt.Errorf("failed to insert product_category: %w", err)
			}
		}
	}

	fmt.Println("✓ Seeders completed successfully")
	return nil
}

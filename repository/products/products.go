package products

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/techies/orders-api/models"
)

type Repository interface {
	Create(ctx context.Context, p *models.Product) error
	GetByID(ctx context.Context, id string) (*models.Product, error)
	Update(ctx context.Context, p *models.Product) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]models.Product, error)
	UpdateStock(ctx context.Context, id string, delta int) error
}

type Repo struct {
	DB  *sql.DB
	log *log.Logger
}

var _ Repository = (*Repo)(nil)

// Create Product
func (r *Repo) Create(ctx context.Context, p *models.Product) error {
	query := `
        INSERT INTO products (id, name, description, price, stock_quantity, sku, is_active, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	now := time.Now()
	_, err := r.DB.ExecContext(ctx, query,
		p.ID, p.Name, p.Description, p.Price, p.StockQuantity, p.SKU, p.IsActive, now, now,
	)
	if err != nil {
		r.log.Printf("Create product error: %v", err)
		return err
	}
	p.CreatedAt = now
	p.UpdatedAt = now
	return nil
}

// GetByID Product
func (r *Repo) GetByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
        SELECT id, name, description, price, stock_quantity, sku, is_active, created_at, updated_at, deleted_at
        FROM products WHERE id = $1 AND deleted_at IS NULL
    `
	var p models.Product
	err := r.DB.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity,
		&p.SKU, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.log.Printf("GetByID error: %v", err)
		return nil, err
	}
	return &p, nil
}

// UpdateStock update product in stock
func (r *Repo) UpdateStock(ctx context.Context, id string, delta int) error {
	query := `
        UPDATE products
        SET stock_quantity = stock_quantity + $1,
            updated_at = $2
        WHERE id = $3 AND deleted_at IS NULL AND is_active = true
    `
	now := time.Now()
	result, err := r.DB.ExecContext(ctx, query, delta, now, id)
	if err != nil {
		r.log.Printf("UpdateStock error: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// List getting the list of products
func (r *Repo) List(ctx context.Context, offset, limit int) ([]models.Product, error) {
	query := `
        SELECT id, name, description, price, stock_quantity, sku, is_active,
               created_at, updated_at, deleted_at
        FROM products
        WHERE deleted_at IS NULL
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `

	rows, err := r.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		r.log.Printf("List error: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			r.log.Printf("List rows close error: %v", err)
		}
	}(rows)
	var products []models.Product
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
			&p.ID, &p.Name, &p.Description, &p.Price, &p.StockQuantity,
			&p.SKU, &p.IsActive, &p.CreatedAt, &p.UpdatedAt, &p.DeletedAt,
		)
		if err != nil {
			r.log.Printf("Scan product error: %v", err)
			return nil, err
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

// Update product
func (r *Repo) Update(ctx context.Context, p *models.Product) error {
	query := `  UPDATE products
        SET name = $1, description = $2, price = $3, stock_quantity = $4,
            sku = $5, is_active = $6, updated_at = $7
        WHERE id = $8 AND deleted_at IS NULL`
	now := time.Now()
	result, err := r.DB.ExecContext(ctx, query, p.Name, p.Description, p.Price, p.StockQuantity,
		p.SKU, p.IsActive, now, p.ID)
	if err != nil {
		r.log.Printf("Update product error: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		r.log.Printf("Update product rows affected error: %v", err)
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	p.UpdatedAt = now
	return nil
}

// Delete Product
func (r *Repo) Delete(ctx context.Context, id string) error {
	query := `UPDATE products SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	now := time.Now()
	result, err := r.DB.ExecContext(ctx, query, now, id)
	if err != nil {
		r.log.Printf("Delete product error: %v", err)
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows // product not found or already deleted
	}
	return nil
}

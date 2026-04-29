package models

import "time"

type Product struct {
	ID            string     `db:"id" json:"id"`
	Name          string     `db:"name" json:"name"`
	Description   string     `db:"description" json:"description"`
	Price         float64    `db:"price" json:"price"`
	StockQuantity int        `db:"stock_quantity" json:"stock_quantity"`
	SKU           *string    `db:"sku" json:"sku,omitempty"`
	IsActive      bool       `db:"is_active" json:"is_active"`
	CreatedAt     time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

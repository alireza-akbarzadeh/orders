package models

import (
	"time"
)

type Order struct {
	ID             string    `db:"id" json:"id"`
	UserID         string    `db:"user_id" json:"user_id"`
	ProductID      *string   `db:"product_id" json:"product_id"`
	TrackingNumber string    `db:"tracking_number" json:"tracking_number"`
	ProductName    string    `db:"product_name" json:"product_name"`
	Quantity       int       `db:"quantity" json:"quantity"`
	UnitPrice      float64   `db:"unit_price" json:"unit_price"`
	TotalPrice     float64   `db:"total_price" json:"total_price"`
	Status         string    `db:"status" json:"status"`
	CanceledAt     time.Time `db:"canceled_at" json:"canceled_at"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

package Orders

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/techies/orders-api/models"
)

type Repo struct {
	DB  *sql.DB
	log *log.Logger
}

// List Order list
func (h *Repo) List(ctx context.Context) ([]models.Order, error) {
	rows, err := h.DB.QueryContext(ctx, `
		SELECT
			id,
			user_id,
			product_name,
			quantity,
			unit_price,
			total_price,
			status,
			created_at,
			updated_at
		FROM orders
		ORDER BY id DESC
	`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			h.log.Printf("failed to close rows: %v", err)
		}
	}()

	orders := make([]models.Order, 0)

	for rows.Next() {
		var b models.Order
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.ProductName,
			&b.Quantity,
			&b.UnitPrice,
			&b.TotalPrice,
			&b.Status,
			&b.CreatedAt,
			&b.UpdatedAt,
		); err != nil {
			return nil, err
		}

		orders = append(orders, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// Create order
func (h *Repo) Create(ctx context.Context, o *models.Order) error {
	query := `
        INSERT INTO orders (
            id, user_id, tracking_number, product_name, quantity,
            unit_price, total_price, status,
            created_at, updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
    `
	_, err := h.DB.ExecContext(ctx, query,
		o.ID,
		o.UserID,
		o.TrackingNumber,
		o.ProductName,
		o.Quantity,
		o.UnitPrice,
		o.TotalPrice,
		o.Status,
		o.CreatedAt,
		o.UpdatedAt,
	)
	return err
}

// Update modifies an existing order. Only specific fields are allowed to change.
func (h *Repo) Update(ctx context.Context, order *models.Order) error {

	query := `
		UPDATE orders SET
			product_name = ?,
			quantity = ?,
			unit_price = ?,
			total_price = ?,
			status = ?,
			updated_at = ?
		WHERE id = ?
	`
	_, err := h.DB.ExecContext(ctx, query,
		order.ProductName,
		order.Quantity,
		order.UnitPrice,
		order.TotalPrice,
		order.Status,
		time.Now(),
		order.ID,
	)
	if err != nil {
		h.log.Printf("failed to update order %s: %v", order.ID, err)
	}
	return err
}

// Delete removes an order by its ID.
func (h *Repo) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = ?`
	_, err := h.DB.ExecContext(ctx, query, id)
	if err != nil {
		h.log.Printf("failed to delete order %s: %v", id, err)
	}
	return err
}

// ListByUser returns all orders placed by a specific user.
func (h *Repo) ListByUser(ctx context.Context, userID string) ([]models.Order, error) {
	query := `
		SELECT
			id, user_id, tracking_number, product_name, quantity,
			unit_price, total_price, status,
			created_at, updated_at
		FROM orders
		WHERE user_id = ?
		ORDER BY created_at DESC
	`

	rows, err := h.DB.QueryContext(ctx, query, userID)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			h.log.Printf("failed to close rows: %v", err)
		}
	}()

	orders := make([]models.Order, 0)
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.ID,
			&o.UserID,
			&o.TrackingNumber,
			&o.ProductName,
			&o.Quantity,
			&o.UnitPrice,
			&o.TotalPrice,
			&o.Status,
			&o.CreatedAt,
			&o.UpdatedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orders, nil
}

// CancelOrder - sets status to "cancelled" and records cancelled_at timestamp.
func (h *Repo) CancelOrder(ctx context.Context, id string) error {
	now := time.Now()
	query := `
        UPDATE orders SET
            status = 'cancelled',
            cancelled_at = ?,
            updated_at = ?
        WHERE id = ? AND status != 'cancelled'
    `
	result, err := h.DB.ExecContext(ctx, query, now, now, id)
	if err != nil {
		h.log.Printf("failed to cancel order %s: %v", id, err)
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

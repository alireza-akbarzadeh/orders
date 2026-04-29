package users

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/techies/orders-api/models"
)

type Repo struct {
	DB  *sql.DB
	log *log.Logger
}

// List Get list of users
func (r *Repo) List(ctx context.Context) ([]models.User, error) {
	rows, err := r.DB.QueryContext(ctx, `
	SELECT
        id,
        username,
        email,
        created_at,
        updated_at,
        is_active,
        role
    FROM users
    ORDER BY created_at DESC
	`)

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			r.log.Printf("failed tp close rows: %v", err)
		}
	}()
	users := make([]models.User, 0)
	for rows.Next() {
		var u models.User
		if err := rows.Scan(
			&u.ID,
			&u.Username,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
			&u.IsActive,
			&u.Role,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

// GetById returns a single user by id

func (r *Repo) GetById(ctx context.Context, id string) (*models.User, error) {

	var u models.User

	row := r.DB.QueryRowContext(ctx, `
SELECT id, username, email, created_at, updated_at, is_active, role
from users WHERE id = $1

`, id)

	err := row.Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.IsActive,
		&u.Role,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

// Update updates an existing user (all fields except password)
func (r *Repo) Update(ctx context.Context, u *models.User) error {
	query := `Update users SET username = ? ,email = ? , updated_at = ? , is_active = ? , role = ? where id = ?`

	_, err := r.DB.ExecContext(ctx, query, u.Username, u.Email, time.Now().UTC(), u.IsActive, u.Role, u.ID)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes a user by ID (optional)
func (r *Repo) Delete(ctx context.Context, id string) error {
	result, err := r.DB.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// UpdatePassword updates only the password hash
func (r *Repo) UpdatePassword(ctx context.Context, id, password string) error {
	query := `UPDATE users SET hash_Password = ? ,updated_at =? where id = ?`
	_, err := r.DB.ExecContext(ctx, query, password, time.Now().UTC(), id)
	return err
}

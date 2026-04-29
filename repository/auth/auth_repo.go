package auth

import (
	"database/sql"
	"time"

	"github.com/techies/orders-api/models"
)

type Repo struct {
	DB        *sql.DB
	JWTSecret []byte
}

// Register inserts a new user with hashed password
func (r *Repo) Register(user *models.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true
	if user.Role == "" {
		user.Role = "user"
	}
	_, err := r.DB.Exec(
		`INSERT INTO users (id, username, email, hashed_password, created_at, updated_at, is_active, role)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.HashedPassword, user.CreatedAt, user.UpdatedAt, user.IsActive, user.Role,
	)
	return err
}

// FindByEmail returns a user by email
func (r *Repo) FindByEmail(email string) (*models.User, error) {
	row := r.DB.QueryRow(
		`SELECT id, username, email, hashed_password, created_at, updated_at, is_active, role FROM users WHERE email = ?`,
		email,
	)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUserName FindByEmail returns a user by email
func (r *Repo) FindByUserName(username string) (*models.User, error) {
	row := r.DB.QueryRow(
		`SELECT id, username, email, hashed_password, created_at, updated_at, is_active, role FROM users WHERE username = ?`,
		username,
	)
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

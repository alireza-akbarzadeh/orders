package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/techies/orders-api/helpers"
	"github.com/techies/orders-api/models"
	"github.com/techies/orders-api/repository/auth"
)

type Auth struct {
	Repo      *auth.Repo
	JWTSecret []byte
}

const (
	Identifier_Required = "identifier and password are required"
)

func (a *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Identifier string `json:"identifier"`
		Password   string `json:"password"`
	}
	if err := helpers.Read(r, &req); err != nil {
		helpers.WriteJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	req.Identifier = strings.TrimSpace(req.Identifier)
	req.Password = strings.TrimSpace(req.Password)
	if req.Identifier == "" || req.Password == "" {
		helpers.WriteJSONError(w, Identifier_Required, http.StatusBadRequest)
		return
	}

	user, err := a.Repo.FindByEmail(req.Identifier)
	if err != nil || user == nil {
		user, err = a.Repo.FindByUserName(req.Identifier)
		if err != nil || user == nil {
			helpers.WriteJSONError(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
	}

	// Validate password (assuming helper exists)
	if !helpers.CheckPasswordHash(req.Password, user.HashedPassword) {
		helpers.WriteJSONError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	token, err := helpers.GenerateJWT(user, a.JWTSecret)
	if err != nil {
		helpers.WriteJSONError(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	helpers.WriteToJSON(w, http.StatusOK, map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (a *Auth) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := helpers.Read(r, &req); err != nil {
		helpers.WriteJSONError(w, "Invalid request", http.StatusBadRequest)
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Email == "" || req.Password == "" {
		helpers.WriteJSONError(w, "All fields are required", http.StatusBadRequest)
		return
	}
	// Check if user already exists
	if user, _ := a.Repo.FindByEmail(req.Email); user != nil {
		helpers.WriteJSONError(w, "Identifier already registered", http.StatusConflict)
		return
	}
	// Hash password using helper
	hash, err := helpers.HashPassword(req.Password)
	if err != nil {
		helpers.WriteJSONError(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user := &models.User{
		ID:             uuid.NewString(),
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hash,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		IsActive:       true,
		Role:           "user",
	}
	if err := a.Repo.Register(user); err != nil {
		helpers.WriteJSONError(w, "Failed to register user", http.StatusInternalServerError)
		return
	}
	helpers.WriteToJSON(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

// FindByEmailHandler returns user info by email (for admin/debug, not for public use)
func (a *Auth) FindByEmailHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email"`
	}
	if err := helpers.Read(r, &req); err != nil || strings.TrimSpace(req.Email) == "" {
		helpers.WriteJSONError(w, "Invalid email", http.StatusBadRequest)
		return
	}
	user, err := a.Repo.FindByEmail(req.Email)
	if err != nil || user == nil {
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}
	// Do not expose hashed password
	user.HashedPassword = ""
	helpers.WriteToJSON(w, http.StatusOK, user)
}

// FindByUserNameHandler returns user info by username (for admin/debug, not for public use)
func (a *Auth) FindByUserNameHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
	}
	if err := helpers.Read(r, &req); err != nil || strings.TrimSpace(req.Username) == "" {
		helpers.WriteJSONError(w, "Invalid username", http.StatusBadRequest)
		return
	}
	user, err := a.Repo.FindByUserName(req.Username)
	if err != nil || user == nil {
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}
	// Do not expose hashed password
	user.HashedPassword = ""
	helpers.WriteToJSON(w, http.StatusOK, user)
}

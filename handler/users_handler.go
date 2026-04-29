package handler

import (
	"net/http"

	"github.com/techies/orders-api/helpers"
	"github.com/techies/orders-api/models"
	"github.com/techies/orders-api/repository/users"
)

type UsersHandler struct {
	Repo *users.Repo
}

// GetUserList List all users
func (h *UsersHandler) GetUserList(w http.ResponseWriter, r *http.Request) {
	list, err := h.Repo.List(r.Context())
	if err != nil {
		helpers.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if list == nil {
		list = make([]models.User, 0)
	}

	helpers.WriteToJSON(w, http.StatusOK, list)
}

// GetUserById Get Single user by ID
func (h *UsersHandler) GetUserById(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.GetUrlParam(w, r, "id")
	if !ok {
		return
	}
	user, err := h.Repo.GetById(r.Context(), id)
	if err != nil {
		helpers.WriteJSONError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if user == nil {
		helpers.WriteJSONError(w, "User not found.", http.StatusNotFound)
		return
	}
	helpers.WriteToJSON(w, http.StatusOK, user)
}

// UpdateUser Update an existing user
func (h *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.GetUrlParam(w, r, "id")
	if !ok {
		return
	}

	// 2. Fetch existing user
	existing, err := h.Repo.GetById(r.Context(), id)
	if err != nil {
		helpers.WriteJSONError(w, "Failed to fetch user: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if existing == nil {
		helpers.WriteJSONError(w, "User not found", http.StatusNotFound)
		return
	}

	var updates struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		IsActive *bool  `json:"is_active"`
		Role     string `json:"role"`
	}

	if err := helpers.DecodeJSON(r, &updates); err != nil {
		helpers.WriteJSONError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if updates.Username != "" {
		existing.Username = updates.Username
	}
	if updates.Email != "" {
		existing.Email = updates.Email
	}
	if updates.IsActive != nil {
		existing.IsActive = *updates.IsActive
	}
	if updates.Role != "" {
		existing.Role = updates.Role
	}
	if err := h.Repo.Update(r.Context(), existing); err != nil {
		helpers.WriteJSONError(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Return updated user (exclude hashed password for safety)
	existing.HashedPassword = ""
	helpers.WriteToJSON(w, http.StatusOK, existing)
}

// Delete User

func (h *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := helpers.GetUrlParam(w, r, "id")
	if !ok {
		return
	}
	if err := h.Repo.Delete(r.Context(), id); err != nil {
		helpers.WriteJSONError(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
	}
	helpers.WriteToJSON(w, http.StatusOK, true)
}

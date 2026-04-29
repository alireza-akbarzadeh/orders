package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/techies/orders-api/helpers"
	"github.com/techies/orders-api/middleware"
	"github.com/techies/orders-api/models"
	Orders "github.com/techies/orders-api/repository/orders"
)

type Order struct {
	Repo *Orders.Repo
}

// ListOrders returns orders for the authenticated user (role‑ignorant)
func (h *Order) ListOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		helpers.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var orders []models.Order
	var err error

	if claims.Role == "admin" {
		orders, err = h.Repo.List(r.Context())
	} else {
		orders, err = h.Repo.ListByUser(r.Context(), claims.UserID)
	}

	if err != nil {
		helpers.WriteJSONError(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}

	if orders == nil {
		orders = []models.Order{}
	}
	helpers.WriteToJSON(w, http.StatusOK, orders)
}

type CreateOrderRequest struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Status      string  `json:"status"`
}

type UpdateOrderRequest struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	TotalPrice  float64 `json:"total_price"`
	Status      string  `json:"status"`
}

// CreateOrder handle user creating order
func (h *Order) CreateOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		helpers.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var req CreateOrderRequest
	if err := helpers.Read(r, &req); err != nil {
		helpers.WriteJSONError(w, "Failed to read request", http.StatusBadRequest)
		return
	}
	// validation
	if req.ProductName == "" || req.Quantity <= 0 || req.UnitPrice <= 0 {
		helpers.WriteJSONError(w, "product_name, quantity (>0) and unit_price (>0) are required", http.StatusBadRequest)
		return
	}
	if req.Status == "" {
		req.Status = "pending"
	}

	totalPrice := float64(req.Quantity) * req.UnitPrice
	now := time.Now()

	order := &models.Order{
		ID:          uuid.NewString(),
		UserID:      claims.UserID,
		ProductName: req.ProductName,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		TotalPrice:  totalPrice,
		Status:      req.Status,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.Repo.Create(r.Context(), order); err != nil {
		helpers.WriteJSONError(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	helpers.WriteToJSON(w, http.StatusCreated, order)
}

// UpdateOrder handle admin updating the order
func (h *Order) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)

	if !ok {
		helpers.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
	}

	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.WriteJSONError(w, "Missing order id", http.StatusBadRequest)
		return
	}
	if claims.Role != "admin" && claims.Role != "manager" {
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}
	var req UpdateOrderRequest

	if err := helpers.Read(r, &req); err != nil {
		helpers.WriteJSONError(w, "Failed to read request", http.StatusBadRequest)
	}
	if req.ProductName == "" || req.Quantity <= 0 || req.UnitPrice <= 0 {
		helpers.WriteJSONError(w, "product_name, quantity (>0) and unit_price (>0) are required", http.StatusBadRequest)
		return
	}
	if req.TotalPrice <= 0 {
		req.TotalPrice = float64(req.Quantity) * req.UnitPrice
	}
	order := &models.Order{
		ID:          id,
		ProductName: req.ProductName,
		Quantity:    req.Quantity,
		UnitPrice:   req.UnitPrice,
		TotalPrice:  req.TotalPrice,
		Status:      req.Status,
	}

	if err := h.Repo.Update(r.Context(), order); err != nil {
		helpers.WriteJSONError(w, "Failed to update order", http.StatusInternalServerError)
		return
	}

	helpers.WriteToJSON(w, http.StatusOK, map[string]string{"message": "Order updated successfully"})
}

// GetUserOrders  returns orders for the authenticated user (role‑ignorant)
func (h *Order) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		helpers.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
	}
	orders, err := h.Repo.ListByUser(r.Context(), claims.UserID)
	if err != nil {
		helpers.WriteJSONError(w, "Failed to fetch orders", http.StatusInternalServerError)
		return
	}
	if orders == nil {
		orders = []models.Order{}
	}
	helpers.WriteToJSON(w, http.StatusOK, orders)
}

// DeleteOrder delete the order
func (h *Order) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		helpers.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if claims.Role != "admin" {
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.WriteJSONError(w, "Missing order id", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Delete(r.Context(), id); err != nil {
		helpers.WriteJSONError(w, "Failed to delete order", http.StatusInternalServerError)
		return
	}
	helpers.WriteToJSON(w, http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}

func (h *Order) CancelOrder(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r)
	if !ok {
		helpers.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Only admin can cancel an order (or you can allow manager too)
	if claims.Role != "admin" {
		helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		return
	}
	id := chi.URLParam(r, "id")
	if id == "" {
		helpers.WriteJSONError(w, "Missing order id", http.StatusBadRequest)
		return
	}

	if err := h.Repo.CancelOrder(r.Context(), id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			helpers.WriteJSONError(w, "Order not found or already cancelled", http.StatusNotFound)
			return
		}
		helpers.WriteJSONError(w, "Failed to cancel order", http.StatusInternalServerError)
		return
	}
	helpers.WriteToJSON(w, http.StatusOK, map[string]string{"message": "Order cancelled successfully"})
}

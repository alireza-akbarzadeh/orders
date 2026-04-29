package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/crypto/bcrypt"
)

// ParseID extracts and parses the "id" URL parameter from the request.
func ParseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	return strconv.ParseInt(idStr, 10, 64)
}

// WriteToJSON writes the given data as a JSON response with the specified status code.
func WriteToJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSONError sends a JSON error response with the specified status code
func WriteJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(ErrorResponse{Error: message}); err != nil {
		http.Error(w, "Failed to encode JSON error response", http.StatusInternalServerError)
	}
}

// Read decodes JSON from the request body into the provided destination structure.
func Read(r *http.Request, dst any) error {
	decoded := json.NewDecoder(r.Body)
	decoded.DisallowUnknownFields()
	return decoded.Decode(dst)
}

// HashPassword hashes a plain password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a plain password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GetUrlParam(w http.ResponseWriter, r *http.Request, param string) (string, bool) {
	value := chi.URLParam(r, param)
	if value == "" {
		WriteJSONError(w, fmt.Sprintf("Please provide a %s.", param), http.StatusBadRequest)
		return "", false
	}
	return value, true
}

func DecodeJSON(r *http.Request, v interface{}) error {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(r.Body)
	return json.NewDecoder(r.Body).Decode(v)
}

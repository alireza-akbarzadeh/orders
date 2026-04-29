package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/techies/orders-api/helpers"
)

type contextKey string

const UserClaim = contextKey("userClaims")

func Authenticate(jwtSecret []byte) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				helpers.WriteJSONError(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				helpers.WriteJSONError(w, "Invalid authorization header", http.StatusUnauthorized)
				return
			}
			tokenString := parts[1]
			claims := &helpers.UserClaims{}

			parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return jwtSecret, nil
			})

			if err != nil || !parsedToken.Valid {
				helpers.WriteJSONError(w, "Invalid or expired token", http.StatusUnauthorized)
				return
			}

			// Store claims in request context – use the defined key UserClaim
			ctx := context.WithValue(r.Context(), UserClaim, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(allowedRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := GetUserClaims(r)
			if !ok {
				helpers.WriteJSONError(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			for _, role := range allowedRoles {
				if claims.Role == role {
					next.ServeHTTP(w, r)
					return
				}
			}
			helpers.WriteJSONError(w, "Forbidden", http.StatusForbidden)
		})
	}
}

// GetUserClaims Get claims
func GetUserClaims(r *http.Request) (claims *helpers.UserClaims, ok bool) {
	claims, ok = r.Context().Value(UserClaim).(*helpers.UserClaims)
	return claims, ok
}

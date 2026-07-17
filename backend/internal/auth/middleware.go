package auth

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

// Authenticates a user
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {

        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        // Remove the "Bearer" prefix
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")

        // Parse and verify the token
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return jwtSecret, nil
        })

        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims := token.Claims.(jwt.MapClaims)

        userID := int(claims["userId"].(float64))
        role := claims["role"].(string)

        ctx := context.WithValue(r.Context(), "userID", userID)
        ctx = context.WithValue(ctx, "role", role)

        // Continue as authentication is complete
        next(w, r.WithContext(ctx))
    }
}
package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Checks if a user is registered and if email and password matches the database
func Login(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
			return
		}

		ctx := r.Context()

		// Variables where the database-read values are stored
		var id int
		var hash string
		var role string

		err := db.QueryRow(ctx, `
            SELECT id, password_hash, role
            FROM users
            WHERE email = $1
        `, req.Email).Scan(&id, &hash, &role)

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid credentials"})
			return
		}

		// Compares the hashed input with the hashed password in the database
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid credentials"})
			return
		}

		// Generate token upon successful login
		token, err := GenerateToken(id, role)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to generate token"})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

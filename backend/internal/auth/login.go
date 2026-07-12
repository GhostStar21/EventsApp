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
		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
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
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Compares the hashed input with the hashed password in the database
		if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)); err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Generate token upon successful login
		token, err := GenerateToken(id, role)
		if err != nil {
			http.Error(w, "Unable to generate token", http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"token": token,
		})
	}
}

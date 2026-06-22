package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		_, err = db.Exec(ctx, `
            INSERT INTO users (name, email, password_hash)
            VALUES ($1, $2, $3)
        `, req.Name, req.Email, string(hash))

		if err != nil {
			http.Error(w, "User already exists", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

package auth

import (
	"encoding/json"
	"net/http"
	// "EventsApp/internal/consts"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Registers a user, storing their credentials in the database
func RegisterUser(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid JSON"})
			return
		}

		// Hash the password
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Error hashing password"})
			return
		}

		// Store the values into the database
		ctx := r.Context()
		_, err = db.Exec(ctx, `
            INSERT INTO users (name, email, password_hash)
            VALUES ($1, $2, $3)
        `, req.Name, req.Email, string(hash))

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "User already exists"})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
	}
}

// RegisterOrganizer promotes an already authenticated user to organizer status
// only when that user exists in the organizer_member table.
// func RegisterOrganizer(db *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		userID, err := GetUserIDFromContext(r.Context())
// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			json.NewEncoder(w).Encode(map[string]string{"message": "You must be logged in to become an organizer"})
// 			return
// 		}

// 		ctx := r.Context()
// 		var membershipExists bool
// 		err = db.QueryRow(ctx, `
// 			SELECT EXISTS (
// 				SELECT 1
// 				FROM organizer_member
// 				WHERE user_id = $1
// 			)
// 		`, userID).Scan(&membershipExists)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to verify organizer membership"})
// 			return
// 		}
// 		if !membershipExists {
// 			w.WriteHeader(http.StatusForbidden)
// 			json.NewEncoder(w).Encode(map[string]string{"message": "You are not registered as an organizer member"})
// 			return
// 		}

// 		_, err = db.Exec(ctx, `
// 			UPDATE users
// 			SET role = $1
// 			WHERE id = $2
// 		`, consts.RoleOrganizer, userID)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to promote user to organizer"})
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(map[string]string{"message": "User promoted to organizer"})
// 	}
// }

// // DemoteOrganizer downgrades an authenticated organizer user back to USER role.
// func DemoteOrganizer(db *pgxpool.Pool) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Content-Type", "application/json")

// 		userID, err := GetUserIDFromContext(r.Context())
// 		if err != nil {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			json.NewEncoder(w).Encode(map[string]string{"message": "You must be logged in"})
// 			return
// 		}

// 		ctx := r.Context()
// 		_, err = db.Exec(ctx, `
// 			UPDATE users
// 			SET role = $1
// 			WHERE id = $2
// 		`, consts.RoleUser, userID)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to demote organizer"})
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		json.NewEncoder(w).Encode(map[string]string{"message": "User demoted to user"})
// 	}
// }
package auth

import (
	"EventsApp/internal/consts"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"github.com/jackc/pgx/v5/pgconn"
	"regexp"
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

		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,4}$`)
            if !emailRegex.MatchString(req.Email) {
                w.WriteHeader(http.StatusBadRequest)
                json.NewEncoder(w).Encode(map[string]string{"message": "Invalid email format"})
                return
            }

		if len(req.Password) < 8 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Password must be at least 8 characters long"})
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
            if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" { // 23505 is unique_violation
                w.WriteHeader(http.StatusConflict) // 409 Conflict is more appropriate for duplicates
                json.NewEncoder(w).Encode(map[string]string{"message": "User with this email already exists"})
                return
            }
            log.Printf("Error registering user: %v", err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(map[string]string{"message": "An internal server error occurred during registration"})
            return
		}

		


		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
	}
}

type organizerRegistrationPayload struct {
	Name      string `json:"name"`
	OrgNumber int    `json:"orgNumber"`
	Type      string `json:"type"`
}

func parseOrganizerRegistrationPayload(r *http.Request) (organizerRegistrationPayload, bool, error) {
	var payload organizerRegistrationPayload
	if r.Body == nil {
		return payload, false, nil
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		if err == io.EOF {
			return payload, false, nil
		}
		return payload, false, err
	}

	if strings.TrimSpace(payload.Name) == "" && payload.OrgNumber == 0 {
		return payload, false, nil
	}

	payload.Type = normalizeOrganizerType(payload.Type)
	return payload, true, nil
}

func normalizeOrganizerType(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "Personal"
	}
	return value
}

// RegisterOrganizer creates or links an organizer for the current authenticated user.
func PromoteOrganizer(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, err := GetUserIDFromContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "You must be logged in to become an organizer"})
			return
		}

		payload, hasPayload, err := parseOrganizerRegistrationPayload(r)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid organizer payload"})
			return
		}

		ctx := r.Context()
		var organizerID int

		if hasPayload {
			if strings.TrimSpace(payload.Name) == "" || payload.OrgNumber == 0 {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"message": "Organizer name and organization number are required"})
				return
			}

			err = db.QueryRow(ctx, `
				INSERT INTO organizers (name, org_number, type)
				VALUES ($1, $2, $3)
				RETURNING id
			`, payload.Name, payload.OrgNumber, payload.Type).Scan(&organizerID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Unable to create organizer"})
				return
			}

			_, err = db.Exec(ctx, `
				INSERT INTO organizer_member (user_id, organizer_id)
				VALUES ($1, $2)
			`, userID, organizerID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Unable to link organizer to user"})
				return
			}
		} else {
			err = db.QueryRow(ctx, `
				SELECT organizer_id
				FROM organizer_member
				WHERE user_id = $1
				ORDER BY organizer_id
				LIMIT 1
			`, userID).Scan(&organizerID)
			if err == pgx.ErrNoRows {
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(map[string]string{"message": "You are not registered as an organizer member"})
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(map[string]string{"message": "Unable to select organizer"})
				return
			}
		}

		_, err = db.Exec(ctx, `
			UPDATE users
			SET role = $1
			WHERE id = $2
		`, consts.RoleOrganizer, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to promote user to organizer"})
			return
		}

		if err := SetSessionCookie(w, userID, string(consts.RoleOrganizer)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to refresh session"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":     "Organizer created",
			"organizerId": organizerID,
		})
	}
}

// DemoteOrganizer downgrades an authenticated organizer user back to USER role.
func DemoteOrganizer(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, err := GetUserIDFromContext(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"message": "You must be logged in"})
			return
		}

		ctx := r.Context()
		_, err = db.Exec(ctx, `
			UPDATE users
			SET role = $1
			WHERE id = $2
		`, consts.RoleUser, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to demote organizer"})
			return
		}

		if err := SetSessionCookie(w, userID, string(consts.RoleUser)); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Unable to refresh session"})
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "User demoted to user"})
	}
}

package users

import (
	"EventsApp/internal/api"
	"EventsApp/internal/auth"
	"EventsApp/internal/consts"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var db *pgxpool.Pool

// Sets the global database connection pool.
func SetDB(pool *pgxpool.Pool) {
	db = pool
}

// Handles different HTTP methods.
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		postUsers(w, r)
	case http.MethodPut:
		updateUsers(w, r)
	case http.MethodDelete:
		deleteUsers(w, r)
	default:
		http.Error(w, "This method ( "+r.Method+" ) is not supported", http.StatusNotImplemented)
		return
	}
}

// Lists either all users or a single user (admin function).
func getUsers(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.UsersPath)
	if isList {
		listUsers(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	getSingleUser(w, r, id)
}

// Lists all users in the database.
func listUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !auth.IsAdmin(ctx) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if db == nil {
		api.WriteJSON(w, users)
		return
	}
	rows, err := db.Query(ctx, `
		SELECT u.id, u.name, u.email, u.role,
		       EXISTS(SELECT 1 FROM organizer_member om WHERE om.user_id = u.id) AS is_organizer_member,
		       (SELECT om.organizer_id FROM organizer_member om WHERE om.user_id = u.id ORDER BY om.organizer_id LIMIT 1) AS organizer_id
		FROM users u
		ORDER BY u.id
	`)
	if err != nil {
		http.Error(w, "Failed to query users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var out []User
	for rows.Next() {
		var u User
		var roleStr string
		var orgID *int
		if err := rows.Scan(&u.Id, &u.Name, &u.Email, &roleStr, &u.IsOrganizerMember, &orgID); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		u.Role = consts.Role(roleStr)
		u.OrganizerId = orgID
		out = append(out, u)
	}
	api.WriteJSON(w, out)
}

// List a single user based on the id provided.
func getSingleUser(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	currentUserID, _ := auth.GetUserIDFromContext(ctx)
	if !auth.IsAdmin(ctx) && currentUserID != id {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if db == nil {
		for _, u := range users {
			if u.Id == id {
				api.WriteJSON(w, u)
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var u User
	var roleStr string
	var orgID *int
	err := db.QueryRow(ctx, `
		SELECT u.id, u.name, u.email, u.role,
		       EXISTS(SELECT 1 FROM organizer_member om WHERE om.user_id = u.id) AS is_organizer_member,
		       (SELECT om.organizer_id FROM organizer_member om WHERE om.user_id = u.id ORDER BY om.organizer_id LIMIT 1) AS organizer_id
		FROM users u
		WHERE u.id = $1
	`, id).Scan(&u.Id, &u.Name, &u.Email, &roleStr, &u.IsOrganizerMember, &orgID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	u.Role = consts.Role(roleStr)
	u.OrganizerId = orgID
	api.WriteJSON(w, u)
}

// Register a user.
type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Role     string `json:"role,omitempty"`
}

func postUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !auth.IsAdmin(ctx) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	passwordHash := ""
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}
		passwordHash = string(hash)
	}

	u := User{
		Name:         req.Name,
		Email:        req.Email,
		Role:         consts.Role(req.Role),
		PasswordHash: passwordHash,
	}

	if db == nil {
		u.Id = len(users) + 1
		users = append(users, u)
		w.WriteHeader(http.StatusCreated)
		api.WriteJSON(w, u)
		return
	}

	if u.Role == "" {
		u.Role = consts.RoleUser
	}

	var id int
	err := db.QueryRow(ctx, "INSERT INTO users (name, email, password_hash, role) VALUES ($1,$2,$3,$4) RETURNING id", u.Name, u.Email, u.PasswordHash, u.Role).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	u.Id = id
	w.WriteHeader(http.StatusCreated)
	api.WriteJSON(w, u)
}

// Update / make changes to user data.
func updateUsers(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.UsersPath)
	if isList {
		http.Error(w, "User id required for update", http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	currentUserID, _ := auth.GetUserIDFromContext(ctx)
	if !auth.IsAdmin(ctx) && currentUserID != id {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if u.Id == 0 {
		u.Id = id
	}
	if u.Id != id {
		http.Error(w, "User ID mismatch", http.StatusBadRequest)
		return
	}

	if db == nil {
		for i, existing := range users {
			if existing.Id == id {
				users[i].Name = u.Name
				users[i].Email = u.Email
				if auth.IsAdmin(ctx) && u.Role != "" {
					users[i].Role = u.Role
				}
				api.WriteJSON(w, users[i])
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if auth.IsAdmin(ctx) {
		if u.Role == "" {
			var currentRole string
			if err := db.QueryRow(ctx, "SELECT role FROM users WHERE id=$1", id).Scan(&currentRole); err != nil {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			u.Role = consts.Role(currentRole)
		}
		cmdTag, err := db.Exec(ctx, "UPDATE users SET name=$1, email=$2, role=$3 WHERE id=$4", u.Name, u.Email, u.Role, u.Id)
		if err != nil || cmdTag.RowsAffected() == 0 {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	} else {
		cmdTag, err := db.Exec(ctx, "UPDATE users SET name=$1, email=$2 WHERE id=$3", u.Name, u.Email, u.Id)
		if err != nil || cmdTag.RowsAffected() == 0 {
			http.Error(w, "Failed to update user", http.StatusInternalServerError)
			return
		}
	}

	api.WriteJSON(w, u)
}

// Delete all or a single user in the database.
func deleteUsers(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.UsersPath)
	if isList {
		deleteAllUsers(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
		return
	}
	deleteSingleUser(w, r, id)
}

// Delete all users from the database.
func deleteAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if !auth.IsAdmin(ctx) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if db == nil {
		users = nil
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if _, err := db.Exec(ctx, "DELETE FROM users"); err != nil {
		http.Error(w, "Failed to delete users", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Delete a single user from the database.
func deleteSingleUser(w http.ResponseWriter, r *http.Request, id int) {
	if db == nil {
		for i, u := range users {
			if u.Id == id {
				users = append(users[:i], users[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
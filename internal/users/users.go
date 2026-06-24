package users

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/http"
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
	rows, err := db.Query(ctx, "SELECT id, name FROM users ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to query users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// User struct-array to store the database elements read.
	var out []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.Id, &u.Name); err != nil {
			http.Error(w, "Failed to scan user", http.StatusInternalServerError)
			return
		}
		out = append(out, u)
	}
	api.WriteJSON(w, out)
}

// List a single user based on the id provided.
func getSingleUser(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	var u User
	err := db.QueryRow(ctx, "SELECT id, name FROM users WHERE id=$1", id).Scan(&u.Id, &u.Name)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	api.WriteJSON(w, u)
}

// Register a user.
func postUsers(w http.ResponseWriter, r *http.Request) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	var id int
	err := db.QueryRow(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING id", u.Name).Scan(&id)
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

	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if u.Id == 0 {
		u.Id = id
	}
	// Check if the URL ID matches the ID of the element read from the database
	if u.Id != id {
		http.Error(w, "User ID mismatch", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "UPDATE users SET name=$1 WHERE id=$2", u.Name, u.Id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
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
		if _, err := db.Exec(ctx, "DELETE FROM users"); err != nil {
			http.Error(w, "Failed to delete users", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
}

// Delete a single user from the database.
func deleteSingleUser(w http.ResponseWriter, r *http.Request, id int) {
	ctx := r.Context()
	cmdTag, err := db.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	if err != nil || cmdTag.RowsAffected() == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
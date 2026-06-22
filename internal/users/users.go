package users

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func SetDB(pool *pgxpool.Pool) {
	db = pool
}

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

func MeHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uidVal := ctx.Value("userID")
	if uidVal == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var userID int
	switch v := uidVal.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	default:
		http.Error(w, "Invalid user id", http.StatusUnauthorized)
		return
	}

	var u User
	var roleStr string
	err := db.QueryRow(ctx, "SELECT id, name, email, role FROM users WHERE id=$1", userID).Scan(&u.Id, &u.Name, &u.Email, &roleStr)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	u.Role = consts.Role(roleStr)
	api.WriteJSON(w, u)
}

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

func listUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := db.Query(ctx, "SELECT id, name FROM users ORDER BY id")
	if err != nil {
		http.Error(w, "Failed to query users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

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

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	id, isList, err := api.ExtractIDFromRequest(r, consts.UsersPath)
	if isList {
		ctx := r.Context()
		if _, err := db.Exec(ctx, "DELETE FROM users"); err != nil {
			http.Error(w, "Failed to delete users", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	}
	if err != nil {
		http.Error(w, "Invalid user id", http.StatusBadRequest)
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

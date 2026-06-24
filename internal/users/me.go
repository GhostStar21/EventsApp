package users

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"net/http"
)


// Function that ensures that a user only can access its own account.
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

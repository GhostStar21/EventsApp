package users

import (
	"EventsApp/internal/api"
	"EventsApp/internal/consts"
	"net/http"
	"log"
	"github.com/jackc/pgx/v5"
)

// Function that ensures that a user only can access its own account.
func MeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
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
		var orgID *int
		err := db.QueryRow(ctx, `
			SELECT u.id, u.name, u.email, u.role,
				EXISTS(SELECT 1 FROM organizer_member om WHERE om.user_id = u.id) AS is_organizer_member,
				(SELECT om.organizer_id FROM organizer_member om WHERE om.user_id = u.id ORDER BY om.organizer_id LIMIT 1) AS organizer_id
			FROM users u
			WHERE u.id = $1
		`, userID).Scan(&u.Id, &u.Name, &u.Email, &roleStr, &u.IsOrganizerMember, &orgID)
		if err != nil {
			if err == pgx.ErrNoRows {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			log.Printf("MeHandler: failed to load user %d: %v", userID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		
		u.Role = consts.Role(roleStr)
		u.OrganizerId = orgID
		api.WriteJSON(w, u)
	} 
	http.Error(w, "Invalid HTTP Method sent", http.StatusMethodNotAllowed)
}

package auth

import (
	"encoding/json"
	"net/http"
	"time"
)

// LogoutUser clears the authentication cookie so the browser no longer sends a valid session.
func LogoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    "",
			Path:     "/",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteNoneMode,
			Expires:  time.Unix(0, 0),
			MaxAge:   -1,
		})
		json.NewEncoder(w).Encode(map[string]string{"message": "Logged out"})
	}
}

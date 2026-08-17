package auth

import (
	"encoding/json"
	"net/http"
)

// CSRFTokenResponse is the response structure for CSRF token requests
type CSRFTokenResponse struct {
	Token string `json:"token"`
}

// GetCSRFToken is an endpoint that provides CSRF tokens to the frontend
// This endpoint should be called by the frontend before making state-changing requests
func GetCSRFToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Generate a new CSRF token
		token, err := GenerateCSRFToken()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"message": "Failed to generate CSRF token"})
			return
		}

		// Store the token
		StoreCSRFToken(token)

		// Set the token in a cookie for JavaScript to read
		SetCSRFCookie(w, token)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(CSRFTokenResponse{Token: token})
	}
}

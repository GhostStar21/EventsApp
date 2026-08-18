package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	// CSRF token cookie name
	CSRFCookieName = "XSRF-TOKEN"
	// CSRF token header name
	CSRFHeaderName = "X-CSRF-TOKEN"
	// Token expiration time (24 hours)
	TokenExpiration = 24 * time.Hour
)

// CSRFToken represents a CSRF token with metadata
type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

// CSRFStore manages CSRF tokens for sessions
type CSRFStore struct {
	mu     sync.RWMutex
	tokens map[string]CSRFToken
}

var csrfStore = &CSRFStore{
	tokens: make(map[string]CSRFToken),
}

// GenerateCSRFToken creates a new random CSRF token
func GenerateCSRFToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", fmt.Errorf("failed to generate CSRF token: %v", err)
	}
	return hex.EncodeToString(token), nil
}

// StoreCSRFToken stores a CSRF token for a session
func StoreCSRFToken(token string) {
	csrfStore.mu.Lock()
	defer csrfStore.mu.Unlock()

	csrfStore.tokens[token] = CSRFToken{
		Token:     token,
		ExpiresAt: time.Now().Add(TokenExpiration),
	}
}

// ValidateCSRFToken checks if a CSRF token is valid and not expired
func ValidateCSRFToken(token string) bool {
	csrfStore.mu.RLock()
	defer csrfStore.mu.RUnlock()

	csrfToken, exists := csrfStore.tokens[token]
	if !exists {
		return false
	}

	// Check if token has expired
	if time.Now().After(csrfToken.ExpiresAt) {
		return false
	}

	return true
}

// InvalidateCSRFToken removes a CSRF token from the store
func InvalidateCSRFToken(token string) {
	csrfStore.mu.Lock()
	defer csrfStore.mu.Unlock()

	delete(csrfStore.tokens, token)
}

// CSRFMiddleware validates CSRF tokens on state-changing requests
// This middleware should be applied to POST, PUT, DELETE, PATCH endpoints
func CSRFMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only validate for state-changing methods
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == http.MethodPatch {
			// Get CSRF token from header
			token := r.Header.Get(CSRFHeaderName)
			if token == "" {
				// Also try to get from cookie as fallback
				cookie, err := r.Cookie(CSRFCookieName)
				if err == nil {
					token = cookie.Value
				}
			}

			// Validate token
			if token == "" || !ValidateCSRFToken(token) {
				http.Error(w, "Invalid or missing CSRF token", http.StatusForbidden)
				return
			}
		}

		// Continue to the next handler
		next(w, r)
	}
}

// SetCSRFCookie sets the CSRF token as a secure HTTP-only cookie
func SetCSRFCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     CSRFCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true, 
		Secure:   true, 
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(TokenExpiration.Seconds()),
	}
	http.SetCookie(w, cookie)
}

// CleanupExpiredTokens removes expired tokens from the store
// This should be called periodically (e.g., in a background goroutine)
func CleanupExpiredTokens() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		csrfStore.mu.Lock()
		now := time.Now()
		for token, csrfToken := range csrfStore.tokens {
			if now.After(csrfToken.ExpiresAt) {
				delete(csrfStore.tokens, token)
			}
		}
		csrfStore.mu.Unlock()
	}
}

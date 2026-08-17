package auth

import (
	"net/http"
	"os"
	"time"
	"log"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func Initialize() {
    secret := os.Getenv("SECRET_KEY")
    if secret == "" {
        log.Fatal("FATAL ERROR: SECRET_KEY environment variable is not set. Cannot generate JWT tokens.")
		return
    }
    jwtSecret = []byte(secret)
    if len(jwtSecret) < 32 { // HS256 recommended minimum length for good security
        log.Fatal("FATAL ERROR: SECRET_KEY is too short. It should be at least 32 bytes for HS256.")
		return
	}
}

// Generates a signed JWT token for the specified user.
func GenerateToken(userID int, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func SetSessionCookie(w http.ResponseWriter, userID int, role string) error {
	token, err := GenerateToken(userID, role)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return nil
}

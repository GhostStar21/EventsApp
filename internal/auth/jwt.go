package auth

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

// Generates a signed JWT token for the specified user.
func GenerateToken(userID int, role string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "userId": userID,
        "role":   role,
        "exp":    time.Now().Add(time.Hour * 24).Unix(),
    })

    return token.SignedString(jwtSecret)
}
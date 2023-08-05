package jwt

import (
	"auth/internal/pkg/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// todo never save secrets in your code
var sampleSecret = "sample-secret"

func NewJWTForUser(user models.User) (string, error) {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":       time.Now().Add(1 * time.Hour),
		"user_id":   user.ID,
		"user_role": user.Role,
	})
	jwtToken, err := tkn.SignedString([]byte(sampleSecret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

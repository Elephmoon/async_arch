package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// todo never save secrets in your code
var sampleSecret = "sample-secret"

type Result struct {
	Valid  bool
	Role   string
	UserID uuid.UUID
}

type Claims struct {
	Role   string    `json:"user_role"`
	UserID string    `json:"user_id"`
	Exp    time.Time `json:"exp"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString string) (Result, error) {
	res := Result{}
	claims := Claims{}

	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(sampleSecret), nil
	})
	if err != nil {
		return res, err
	}
	if !token.Valid || claims.Exp.Before(time.Now()) {
		return res, nil
	}
	res.UserID, err = uuid.Parse(claims.UserID)
	if err != nil {
		return res, err
	}
	res.Valid = true
	res.Role = claims.Role

	return res, nil
}

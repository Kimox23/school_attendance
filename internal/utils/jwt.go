package utils

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type JWTUtil struct {
	secret     string
	expiration time.Duration
}

func NewJWTUtil(secret string, expiration time.Duration) *JWTUtil {
	return &JWTUtil{
		secret:     secret,
		expiration: expiration,
	}
}

func (j *JWTUtil) GenerateToken(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(j.expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWTUtil) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.ErrUnauthorized
		}
		return []byte(j.secret), nil
	})
}

func (j *JWTUtil) GetClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := j.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fiber.ErrUnauthorized
	}

	return claims, nil
}

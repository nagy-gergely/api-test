package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nagy-gergely/api-test/types"
)

func GenerateToken(user *types.User) (string, error) {
	claims := jwt.MapClaims{
		"id":         user.ID,
		"fist_name":  user.FirstName,
		"last_name":  user.LastName,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("secret"))
}

func ValidateToken(tokenString string) (*types.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	user := &types.User{
		ID:        claims["id"].(int),
		FirstName: claims["first_name"].(string),
		LastName:  claims["last_name"].(string),
		Email:     claims["email"].(string),
		CreatedAt: claims["created_at"].(time.Time),
	}

	return user, nil
}

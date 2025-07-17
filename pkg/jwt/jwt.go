package jwt

import (
	"errors"
	"fmt"
	"os"

	"github.com/dgrijalva/jwt-go"
)

var ErrorEmptyToker error = errors.New("error emprty token")

//nolint:gochecknoglobals // тише тише тише
var jwtSecretKey = []byte(os.Getenv("jwtSecretKey"))

func GenerateTokenAccess(user string) (string, error) {
	payload := jwt.MapClaims{
		"user": user,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	t, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return t, fmt.Errorf("error create token: %w", err)
	}
	return t, nil
}

type UserClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func GetInfoFromToken(token string) (*UserClaims, error) {
	if token == "" {
		return nil, ErrorEmptyToker
	}
	t, err := jwt.ParseWithClaims(token, &UserClaims{}, func(*jwt.Token) (interface{}, error) {
		return jwtSecretKey, nil
	})

	//nolint:revive,wrapcheck // взял из примера
	if claims, ok := t.Claims.(*UserClaims); ok && t.Valid {
		return claims, nil
	} else {
		return claims, err
	}
}

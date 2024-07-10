package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Id   uint32 `json:"uid"`
	Name string
	jwt.RegisteredClaims
}

var secret_key = []byte("Secret Key of UMeet")

func GenerateToken(id uint32, name string) (string, error) {
	claims := UserClaims{
		Id:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "UMeet",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret_key)

}

func ParseToken(token string) (UserClaims, error) {
	var data UserClaims
	valid, err := jwt.ParseWithClaims(token, &data, func(t *jwt.Token) (interface{}, error) {
		return secret_key, nil
	})
	if err != nil || !valid.Valid {
		err = errors.New("invalid token")
	}
	return data, err

}

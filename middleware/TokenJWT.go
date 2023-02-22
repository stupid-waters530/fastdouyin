package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type Claims struct {
	UserId int64
	jwt.StandardClaims
}

var jwtKey = []byte("shliang")

func CreateToken(id int64) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour) //终止时间，7天后
	claims := &Claims{
		UserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "shliang",
			Subject:   "user token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, err
}

func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return token, claims, err
}

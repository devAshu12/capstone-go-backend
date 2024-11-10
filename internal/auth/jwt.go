package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtAccessKey = []byte(os.Getenv("JWT_SECRET"))
var jwtRefreshKey = []byte(os.Getenv("JWT_SECRET_REFRESH"))

func GenerateToken(userID string) (string, string, error) {
	access_claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	refresh_claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * 24 * time.Hour)),
	}

	access_token, err1 := jwt.NewWithClaims(jwt.SigningMethodES256, access_claims).SignedString(jwtAccessKey)
	refresh_token, err2 := jwt.NewWithClaims(jwt.SigningMethodES256, refresh_claims).SignedString(jwtRefreshKey)

	if err1 != nil {
		return "", "", err1
	}
	if err2 != nil {
		return "", "", err2
	}

	return access_token, refresh_token, nil
}

func ValidateAccessToken(tokenString string, isAccessToken bool) (*jwt.RegisteredClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if isAccessToken {
			return jwtAccessKey, nil
		}
		return jwtRefreshKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

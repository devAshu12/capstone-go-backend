package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtAccessKey = []byte(os.Getenv("JWT_SECRET"))
var jwtRefreshKey = []byte(os.Getenv("JWT_SECRET_REFRESH"))

// GenerateToken creates an access and refresh token for a user with the specified userID.
func GenerateToken(userID string) (string, string, error) {
	accessClaims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
	}

	refreshClaims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * 24 * time.Hour)),
	}

	// Use HS256 for HMAC-based symmetric signing
	accessToken, err1 := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(jwtAccessKey)
	refreshToken, err2 := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(jwtRefreshKey)

	if err1 != nil {
		return "", "", err1
	}
	if err2 != nil {
		return "", "", err2
	}

	return accessToken, refreshToken, nil
}

// ValidateAccessToken validates the provided token string as either an access or refresh token.
func ValidateAccessToken(tokenString string, isAccessToken bool) (*jwt.RegisteredClaims, error) {
	// Set the appropriate key based on the token type
	var signingKey []byte
	if isAccessToken {
		signingKey = jwtAccessKey
	} else {
		signingKey = jwtRefreshKey
	}

	// Parse and validate the token with claims
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method is HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return signingKey, nil
	})

	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

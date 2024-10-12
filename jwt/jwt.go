package jwt

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims struct to be used for both tokens
type Claims struct {
	Email  string `json:"email"`
	RoleID *int   `json:"roleId"`
	UserID int    `json:"userId"`
	jwt.RegisteredClaims
}

// CreateAccessToken generates a new access token with a short expiration time
func CreateAccessToken(email string, roleID int, userID int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Set to 15 minutes
	claims := &Claims{
		Email:  email,
		RoleID: &roleID,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	accessSecret, ok := os.LookupEnv("JWT_ACCESS_TOKEN_SECRET")
	if !ok {
		return "", errors.New("access secret not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(accessSecret))
}

// CreateRefreshToken generates a refresh token with a longer expiration time
func CreateRefreshToken(userID int) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour) // Set to 30 days
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	refreshSecret, ok := os.LookupEnv("JWT_REFRESH_TOKEN_SECRET")
	if !ok {
		return "", errors.New("refresh secret not found")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(refreshSecret))
}

// ValidateToken validates the given JWT string and returns the claims if valid
func ValidateToken(tokenStr string, isRefresh bool) (bool, *Claims, error) {
	claims := &Claims{}
	var secret string
	if isRefresh {
		s, ok := os.LookupEnv("JWT_REFRESH_TOKEN_SECRET")
		if !ok {
			return false, nil, errors.New("refresh secret not found")
		}
		secret = s
	} else {
		s, ok := os.LookupEnv("JWT_ACCESS_TOKEN_SECRET")
		if !ok {
			return false, nil, errors.New("access secret not found")
		}
		secret = s
	}

	secretByte := []byte(secret)

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secretByte, nil
	})

	if err != nil || !token.Valid {
		slog.Info("Invalid token", "err", err)
		return false, nil, errors.New("invalid token")
	}

	return true, claims, nil
}

func GetTokenClaims(ctx context.Context) *Claims {
	token := ctx.Value("access_token")
	if token == nil {
		return nil
	}
	jwt, err := jwt.ParseWithClaims(token.(string), &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_ACCESS_TOKEN_SECRET")), nil
	})
	if err != nil {
		return nil
	}
	return jwt.Claims.(*Claims)
}

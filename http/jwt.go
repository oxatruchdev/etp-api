package http

import (
	"errors"
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
func (s *Server) CreateAccessToken(email string, roleID int, userID int) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute) // Set to 15 minutes
	claims := &Claims{
		Email:  email,
		RoleID: &roleID,
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Config.JWTAccessSecret))
}

// CreateRefreshToken generates a refresh token with a longer expiration time
func (s *Server) CreateRefreshToken(userID int) (string, error) {
	expirationTime := time.Now().Add(30 * 24 * time.Hour) // Set to 30 days
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.Config.JWTRefreshSecret))
}

// ValidateToken validates the given JWT string and returns the claims if valid
func (s *Server) ValidateToken(tokenStr string, isRefresh bool) (*Claims, error) {
	claims := &Claims{}
	var secret []byte
	if isRefresh {
		secret = []byte(s.Config.JWTRefreshSecret)
	} else {
		secret = []byte(s.Config.JWTAccessSecret)
	}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

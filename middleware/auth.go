package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Evalua-Tu-Profe/etp-api/jwt"
)

const (
	IsAuthKey      = "is_authenticated"
	AccessTokenKey = "access_token"
	UserIDKey      = "user_id"
	RoleIDKey      = "role_id"
	ClaimsKey      = "claims"
)

func AuthOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("auth optional")
		cookie, err := r.Cookie(AccessTokenKey)
		if err != nil {
			ctx := context.WithValue(r.Context(), IsAuthKey, false)
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
			return
		}

		accessToken := cookie.Value
		isValid, claims, err := jwt.ValidateToken(accessToken, false)
		if err != nil {
			http.SetCookie(w, &http.Cookie{
				Name:     "refresh_token",
				Value:    "",
				Expires:  time.Now(),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "access_token",
				Value:    "",
				Expires:  time.Now(),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			})

			log.Println("entered in middleware is not valid")
			ctx := context.WithValue(r.Context(), IsAuthKey, false)
			req := r.WithContext(ctx)
			next.ServeHTTP(w, req)
			return
		}

		ctx := context.WithValue(r.Context(), IsAuthKey, isValid)

		ctx = context.WithValue(ctx, ClaimsKey, claims)
		ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, RoleIDKey, claims.RoleID)

		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}

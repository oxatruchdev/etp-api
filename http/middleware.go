package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("access_token")
		if err != nil {
			// User not authenticated
			c.Set("isAuthenticated", false)
			return next(c)
		}

		// Validate token
		claims, err := s.ValidateToken(cookie.Value, false)
		if err != nil {
			// Invalid token, treat as unauthenticated
			c.Set("isAuthenticated", false)
			return next(c)
		}

		// User is authenticated, set user info
		c.Set("user", claims.Email)
		c.Set("isAuthenticated", true)
		return next(c)
	}
}

func (s *Server) RequireAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		isAuthenticated := c.Get("isAuthenticated").(bool)

		if !isAuthenticated {
			// Return 401 if not authenticated
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		// User is authenticated, proceed with the request
		return next(c)
	}
}

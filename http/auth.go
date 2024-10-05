package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web/auth"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) registerAuthRoutes() {
	s.Echo.GET("/register", s.register)
	s.Echo.POST("/register", s.createUser)

	s.Echo.GET("/login", s.HandleLogin)
	s.Echo.POST("/login", s.login)
}

func (s *Server) register(c echo.Context) error {
	return Render(c, http.StatusOK, auth.Register(auth.RegisterFormProps{
		Errors: make(map[string]string),
	}))
}

func (s *Server) createUser(c echo.Context) error {
	c.Logger().Info("Registering user")
	var user etp.User
	if err := c.Bind(&user); err != nil {
		slog.Info("Error getting schools", "error", err)
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	userFound, err := s.UserService.GetUserByEmail(c.Request().Context(), user.Email)
	if err != nil {
		return Render(c, http.StatusInternalServerError, auth.RegisterForm(auth.RegisterFormProps{
			Email:    user.Email,
			Password: user.Password,
			Errors:   map[string]string{"message": "Error al registrar usuario, intenta más tarde."},
		}))
	}

	if userFound != nil {
		return Render(c, http.StatusConflict, auth.RegisterForm(auth.RegisterFormProps{
			Email:    user.Email,
			Password: user.Password,
			Errors:   map[string]string{"message": "Este correo ya está en uso, intenta con otro"},
		}))
	}

	if err := s.UserService.RegisterUser(c.Request().Context(), &etp.User{
		Email:    user.Email,
		Password: user.Password,
	}); err != nil {
		c.Logger().Error("error registering user", "error", err, "message", err.Error())
		return Error(c, err)
	}

	return Render(c, http.StatusCreated, auth.SuccessfulRegistration())
}

func (s *Server) HandleLogin(c echo.Context) error {
	return Render(c, http.StatusOK, auth.LoginPage(auth.LoginFormProps{}))
}

func (s *Server) login(c echo.Context) error {
	c.Logger().Info("Logging in user")
	var body etp.User
	if err := c.Bind(&body); err != nil {
		return Error(c, etp.Errorf(etp.EINVALID, "invalid body"))
	}

	foundUser, err := s.UserService.GetUserByEmail(c.Request().Context(), body.Email)
	if err != nil {
		c.Logger().Error("error logging in user ", "error ", err, "message ", err.Error())
		return Error(c, err)
	}

	if foundUser == nil || bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(body.Password)) != nil {
		return Render(c, http.StatusUnauthorized, auth.LoginForm(auth.LoginFormProps{
			Email:    body.Email,
			Password: body.Password,
			Errors: map[string]string{
				"message": "Usuario o contraseña incorrectos",
			},
		}))
	}

	token, err := s.CreateAccessToken(foundUser.Email, *foundUser.RoleID, foundUser.ID)
	if err != nil {
		return Error(c, err)
	}

	refreshToken, err := s.CreateRefreshToken(foundUser.ID)
	if err != nil {
		return Error(c, err)
	}

	c.SetCookie(&http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	slog.Info("tokens generated", "access_token", token, "refresh_token", refreshToken)
	c.SetCookie(&http.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(12 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

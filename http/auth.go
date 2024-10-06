package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web/auth"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAuthRoutes registers the authentication routes
func (s *Server) registerAuthRoutes() {
	// Use standard http.HandleFunc for routing
	http.HandleFunc("GET /register", s.register)
	http.HandleFunc("POST /register", s.createUser)
	http.HandleFunc("GET /login", s.HandleLogin)
	http.HandleFunc("POST /login", s.login)
}

// register renders the registration page
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	Render(w, r, http.StatusOK, auth.Register(auth.RegisterFormProps{
		Errors: make(map[string]string),
	}))
}

// createUser handles the registration process
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	slog.Info("Registering user")

	// Parse the request body
	var user etp.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Info("Error parsing user data", "error", err)
		Error(w, r, etp.Errorf(etp.EINVALID, "invalid body"))
		return
	}

	// Check if the user already exists
	userFound, err := s.UserService.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		Render(w, r, http.StatusInternalServerError, auth.RegisterForm(auth.RegisterFormProps{
			Email:    user.Email,
			Password: user.Password,
			Errors:   map[string]string{"message": "Error al registrar usuario, intenta más tarde."},
		}))
		return
	}

	if userFound != nil {
		Render(w, r, http.StatusConflict, auth.RegisterForm(auth.RegisterFormProps{
			Email:    user.Email,
			Password: user.Password,
			Errors:   map[string]string{"message": "Este correo ya está en uso, intenta con otro"},
		}))
		return
	}

	// Get the student role
	studentRole, err := s.RoleService.GetRoleByName(r.Context(), etp.RoleStudent)
	if err != nil {
		slog.Error("Error getting student role", "error", err)
		Error(w, r, err)
		return
	}

	// Register the user
	if err := s.UserService.RegisterUser(r.Context(), &etp.User{
		Email:    user.Email,
		Password: user.Password,
		RoleID:   &studentRole.ID,
	}); err != nil {
		slog.Error("Error registering user", "error", err)
		Error(w, r, err)
		return
	}

	Render(w, r, http.StatusCreated, auth.SuccessfulRegistration())
}

// HandleLogin renders the login page
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	Render(w, r, http.StatusOK, auth.LoginPage(auth.LoginFormProps{}))
}

// login handles the user login process
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	slog.Info("Logging in user")
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse the request body
	var body etp.User
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		Error(w, r, etp.Errorf(etp.EINVALID, "invalid body"))
		return
	}

	// Find the user by email
	foundUser, err := s.UserService.GetUserByEmail(r.Context(), body.Email)
	if err != nil || foundUser == nil || bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(body.Password)) != nil {
		Render(w, r, http.StatusUnauthorized, auth.LoginForm(auth.LoginFormProps{
			Email:    body.Email,
			Password: body.Password,
			Errors:   map[string]string{"message": "Usuario o contraseña incorrectos"},
		}))
		return
	}

	// Create tokens
	token, err := s.CreateAccessToken(foundUser.Email, *foundUser.RoleID, foundUser.ID)
	if err != nil {
		Error(w, r, err)
		return
	}

	refreshToken, err := s.CreateRefreshToken(foundUser.ID)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(30 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(12 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect or return success
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

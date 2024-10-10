package http

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/Evalua-Tu-Profe/etp-api/cmd/web/auth"
	"github.com/Evalua-Tu-Profe/etp-api/jwt"
	"github.com/Evalua-Tu-Profe/etp-api/middleware"
	"golang.org/x/crypto/bcrypt"
)

// RegisterAuthRoutes registers the authentication routes
func (s *Server) registerAuthRoutes() {
	// Use standard http.HandleFunc for routing
	s.Mux.HandleFunc("GET /register", s.register)
	s.Mux.HandleFunc("POST /register", s.createUser)
	s.Mux.HandleFunc("GET /login", s.HandleLogin)
	s.Mux.HandleFunc("POST /login", s.login)
	s.Mux.HandleFunc("POST /logout", s.logout)
}

// register renders the registration page
func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	slog.Info("Registering user")
	ctx := r.Context()
	isAuth := ctx.Value(middleware.IsAuthKey)

	if isAuth == true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	slog.Info("Is auth", "isAuth", isAuth)
	Render(w, r, http.StatusOK, auth.Register(auth.RegisterFormProps{
		Errors: make(map[string]string),
	}))
}

// createUser handles the registration process
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	slog.Info("Request headers", "headers", r.Header)
	errors := make(map[string]string)
	slog.Info("Registering user")

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}

	if !r.Form.Has("email") {
		errors["email"] = "El correo es obligatorio"
	}

	if !r.Form.Has("password") {
		errors["password"] = "La contraseña es obligatoria"
	}

	var user etp.User

	if len(errors) > 0 {
		Render(w, r, http.StatusBadRequest, auth.RegisterForm(auth.RegisterFormProps{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
			Errors:   errors,
		}))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			slog.Error("Error rendering in post register: ", "error", err)
		}

	} else {
		user = etp.User{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
		}
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
		slog.Info("User already exists")
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

	w.Header().Set("HX-Reswap", "innerHTML")
	Render(w, r, http.StatusCreated, auth.SuccessfulRegistration())
}

// HandleLogin renders the login page
func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	slog.Info("Rendering login page")
	ctx := r.Context()
	isAuth := ctx.Value(middleware.IsAuthKey)

	if isAuth == true {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	Render(w, r, http.StatusOK, auth.LoginPage(auth.LoginFormProps{}))
}

// login handles the user login process
func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var errors map[string]string = make(map[string]string)
	slog.Info("Logging in user")

	// Parse the request body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// make this statement work, idk the exact syntax for this
	if v := r.Form.Get("email"); v == "" || !r.Form.Has("email") {
		errors["email"] = "El correo es obligatorio"
	}

	if !r.Form.Has("password") {
		errors["password"] = "La contraseña es obligatoria"
	}

	slog.Info("Errors", "errors", errors)

	if len(errors) > 0 {
		slog.Info("Errors", "errors greater than 0", errors)
		Render(w, r, http.StatusBadRequest, auth.LoginForm(auth.LoginFormProps{
			Email:    r.Form.Get("email"),
			Password: r.Form.Get("password"),
			Errors:   errors,
		}))
		return
	}

	var body auth.LoginFormProps
	body.Email = r.Form.Get("email")
	body.Password = r.Form.Get("password")

	// Find the user by email
	foundUser, err := s.UserService.GetUserByEmail(r.Context(), body.Email)
	if err != nil || foundUser == nil || bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(body.Password)) != nil {
		slog.Info("User not found")
		Render(w, r, http.StatusUnauthorized, auth.LoginForm(auth.LoginFormProps{
			Email:    body.Email,
			Password: body.Password,
			Errors:   map[string]string{"message": "Usuario o contraseña incorrectos"},
		}))
		return
	}

	// Create tokens
	token, err := jwt.CreateAccessToken(foundUser.Email, *foundUser.RoleID, foundUser.ID)
	if err != nil {
		Error(w, r, err)
		return
	}

	refreshToken, err := jwt.CreateRefreshToken(foundUser.ID)
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
	w.Header().Set("HX-Redirect", "/")
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
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
	w.Header().Set("HX-Redirect", "/")
}

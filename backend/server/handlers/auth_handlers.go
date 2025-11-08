package handlers

import (
	"backend-avanzada/api"
	"backend-avanzada/models"
	"backend-avanzada/repository"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Estructura que representa los claims del token JWT.
// Se define localmente para no depender del paquete server.
type AuthClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// AuthHandler gestiona login y registro de usuarios.
type AuthHandler struct {
	UserRepository repository.UserRepository
	Logger         func(status int, path string, start time.Time)
	HandleError    func(w http.ResponseWriter, statusCode int, path string, cause error)
	JWTSecret      string
}

// Constructor del handler (inyección de dependencias)
func NewAuthHandler(jwtSecret string, ur repository.UserRepository,
	handleError func(w http.ResponseWriter, statusCode int, path string, cause error),
	logger func(status int, path string, start time.Time)) *AuthHandler {

	return &AuthHandler{
		UserRepository: ur,
		JWTSecret:      jwtSecret,
		HandleError:    handleError,
		Logger:         logger,
	}
}

// Registro de nuevo usuario
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req api.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleError(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}
	if req.Email == "" || req.Password == "" || req.Role == "" {
		h.HandleError(w, http.StatusBadRequest, r.URL.Path, errors.New("email, password and role are required"))
		return
	}

	exists, err := h.UserRepository.FindByEmail(req.Email)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if exists != nil {
		h.HandleError(w, http.StatusConflict, r.URL.Path, errors.New("email already registered"))
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	u := &models.User{
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         req.Role,
	}
	if _, err := h.UserRepository.Save(u); err != nil {
		h.HandleError(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Inicio de sesión
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req api.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.HandleError(w, http.StatusBadRequest, r.URL.Path, err)
		return
	}

	u, err := h.UserRepository.FindByEmail(req.Email)
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}
	if u == nil {
		h.HandleError(w, http.StatusUnauthorized, r.URL.Path, errors.New("invalid credentials"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)); err != nil {
		h.HandleError(w, http.StatusUnauthorized, r.URL.Path, errors.New("invalid credentials"))
		return
	}

	now := time.Now()
	claims := &AuthClaims{
		Email: u.Email,
		Role:  u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(2 * time.Hour)),
			Issuer:    "alchemist-system",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		h.HandleError(w, http.StatusInternalServerError, r.URL.Path, err)
		return
	}

	resp := api.AuthResponse{Token: tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

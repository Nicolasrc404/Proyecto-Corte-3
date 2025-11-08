package server

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type userContextKey struct{}

type AuthClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// AuthMiddleware valida el JWT y (opcionalmente) exige uno de los roles dados.
func (s *Server) AuthMiddleware(roles ...string) func(http.Handler) http.Handler {
	// Preparamos un set de roles permitidos (si roles fue provisto)
	roleRequired := map[string]struct{}{}
	for _, r := range roles {
		roleRequired[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				s.HandleError(w, http.StatusUnauthorized, r.URL.Path, errors.New("missing bearer token"))
				return
			}
			tokenString := strings.TrimPrefix(auth, "Bearer ")

			claims := &AuthClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				return []byte(s.jwtSecret), nil
			})
			if err != nil || !token.Valid {
				s.HandleError(w, http.StatusUnauthorized, r.URL.Path, errors.New("invalid token"))
				return
			}

			// Si se especificaron roles, revisamos que el del token esté permitido
			if len(roleRequired) > 0 {
				if _, ok := roleRequired[claims.Role]; !ok {
					s.HandleError(w, http.StatusForbidden, r.URL.Path, errors.New("forbidden"))
					return
				}
			}

			ctx := context.WithValue(r.Context(), userContextKey{}, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Helper para extraer claims si los necesitas en handlers
func GetAuthClaims(r *http.Request) *AuthClaims {
	if v := r.Context().Value(userContextKey{}); v != nil {
		if c, ok := v.(*AuthClaims); ok {
			return c
		}
	}
	return nil
}

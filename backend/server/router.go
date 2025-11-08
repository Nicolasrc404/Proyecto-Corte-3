package server

import (
	"backend-avanzada/server/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) router() http.Handler {
	router := mux.NewRouter()
	router.Use(s.logger.RequestLogger)
	// Auth
	authHandler := handlers.NewAuthHandler(
		s.GetJWTSecret(),
		s.UserRepository,
		s.HandleError,
		s.logger.Info,
	)
	router.HandleFunc("/auth/register", authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)

	// People (GET públicos por ahora, mutaciones protegidas)
	router.HandleFunc("/people", s.HandlePeople).Methods(http.MethodGet)
	router.Handle("/people", s.AuthMiddleware("alchemist", "supervisor")(http.HandlerFunc(s.HandlePeople))).Methods(http.MethodPost)

	router.HandleFunc("/people/{id}", s.HandlePeopleWithId).Methods(http.MethodGet)
	router.Handle("/people/{id}", s.AuthMiddleware("alchemist", "supervisor")(http.HandlerFunc(s.HandlePeopleWithId))).Methods(http.MethodPut, http.MethodDelete)

	// Kills (solo lectura pública; creación requiere supervisor)
	router.HandleFunc("/kills", s.HandleKills).Methods(http.MethodGet)
	router.Handle("/kills/{id}", s.AuthMiddleware("supervisor")(http.HandlerFunc(s.HandleKillsWithId))).Methods(http.MethodPost)

	return router
}

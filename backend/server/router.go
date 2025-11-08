package server

import (
	"backend-avanzada/server/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) router() http.Handler {
	router := mux.NewRouter()
	router.Use(s.logger.RequestLogger)

	// ========== AUTH ==========
	authHandler := handlers.NewAuthHandler(
		s.GetJWTSecret(),
		s.UserRepository,
		s.HandleError,
		s.logger.Info,
	)
	router.HandleFunc("/auth/register", authHandler.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", authHandler.Login).Methods(http.MethodPost)

	// ========== PEOPLE ==========
	// GET públicos; mutaciones protegidas para "alchemist" o "supervisor"
	router.HandleFunc("/people", s.HandlePeople).Methods(http.MethodGet)
	router.Handle(
		"/people",
		s.AuthMiddleware("alchemist", "supervisor")(http.HandlerFunc(s.HandlePeople)),
	).Methods(http.MethodPost)

	router.HandleFunc("/people/{id}", s.HandlePeopleWithId).Methods(http.MethodGet)
	router.Handle(
		"/people/{id}",
		s.AuthMiddleware("alchemist", "supervisor")(http.HandlerFunc(s.HandlePeopleWithId)),
	).Methods(http.MethodPut, http.MethodDelete)

	// ========== KILLS ==========
	// Solo lectura pública; creación requiere "supervisor"
	router.HandleFunc("/kills", s.HandleKills).Methods(http.MethodGet)
	router.Handle(
		"/kills/{id}",
		s.AuthMiddleware("supervisor")(http.HandlerFunc(s.HandleKillsWithId)),
	).Methods(http.MethodPost)

	// ========== ALCHEMISTS ==========
	// Se registran solo si el repo está disponible (tu mismo patrón)
	if s.AlchemistRepository != nil {
		alchHandler := handlers.NewAlchemistHandler(
			s.AlchemistRepository,
			s.PeopleRepository,
			s.HandleError,
			s.logger.Info,
		)
		// Lectura pública
		router.HandleFunc("/alchemists", alchHandler.GetAll).Methods(http.MethodGet)
		router.HandleFunc("/alchemists/{id}", alchHandler.GetByID).Methods(http.MethodGet)

		// Mutaciones protegidas
		router.Handle(
			"/alchemists",
			s.AuthMiddleware("supervisor")(http.HandlerFunc(alchHandler.Create)),
		).Methods(http.MethodPost)
		router.Handle(
			"/alchemists/{id}",
			s.AuthMiddleware("supervisor")(http.HandlerFunc(alchHandler.Edit)),
		).Methods(http.MethodPut)
		router.Handle(
			"/alchemists/{id}",
			s.AuthMiddleware("supervisor")(http.HandlerFunc(alchHandler.Delete)),
		).Methods(http.MethodDelete)
	}

	return router
}

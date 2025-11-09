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

	// ========== ALCHEMISTS ==========
	// Se registran solo si el repo está disponible (tu mismo patrón)
	if s.AlchemistRepository != nil {
		alchHandler := handlers.NewAlchemistHandler(
			s.AlchemistRepository,
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

		// ======== TRANSMUTATIONS ========
		if s.TransmutationRepository != nil {
			transHandler := handlers.NewTransmutationHandler(
				s.TransmutationRepository,
				s.HandleError,
				s.logger.Info,
			)

			router.Handle(
				"/transmutations",
				s.AuthMiddleware("alchemist", "supervisor")(http.HandlerFunc(transHandler.Create)),
			).Methods(http.MethodPost)

			router.HandleFunc("/transmutations", transHandler.GetAll).Methods(http.MethodGet)
			router.HandleFunc("/transmutations/{id}", transHandler.GetByID).Methods(http.MethodGet)

			router.Handle("/transmutations/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(transHandler.Delete)),
			).Methods(http.MethodDelete)

		}

		// ======== MISSIONS ========
		if s.MissionRepository != nil {
			mh := handlers.NewMissionHandler(s.MissionRepository, s.HandleError, s.logger.Info)
			router.HandleFunc("/missions", mh.GetAll).Methods(http.MethodGet)
			router.HandleFunc("/missions/{id}", mh.GetByID).Methods(http.MethodGet)
			router.Handle("/missions",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(mh.Create)),
			).Methods(http.MethodPost)
			router.Handle("/missions/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(mh.Edit)), // <- NUEVO PUT
			).Methods(http.MethodPut)
			router.Handle("/missions/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(mh.Delete)),
			).Methods(http.MethodDelete)
		}

		// ======== TRANSMUTATIONS ========
		if s.TransmutationRepository != nil {
			transHandler := handlers.NewTransmutationHandler(
				s.TransmutationRepository,
				s.HandleError,
				s.logger.Info,
			)

			router.Handle(
				"/transmutations",
				s.AuthMiddleware("alchemist", "supervisor")(http.HandlerFunc(transHandler.Create)),
			).Methods(http.MethodPost)

			router.HandleFunc("/transmutations", transHandler.GetAll).Methods(http.MethodGet)
			router.HandleFunc("/transmutations/{id}", transHandler.GetByID).Methods(http.MethodGet)

			router.Handle(
				"/transmutations/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(transHandler.Edit)), // ✅ nuevo PUT
			).Methods(http.MethodPut)

			router.Handle(
				"/transmutations/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(transHandler.Delete)),
			).Methods(http.MethodDelete)
		}

		// ======== MATERIALS ========
		if s.MaterialRepository != nil {
			matHandler := handlers.NewMaterialHandler(s.MaterialRepository, s.HandleError, s.logger.Info)
			router.HandleFunc("/materials", matHandler.GetAll).Methods(http.MethodGet)
			router.HandleFunc("/materials/{id}", matHandler.GetByID).Methods(http.MethodGet)
			router.Handle("/materials",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(matHandler.Create)),
			).Methods(http.MethodPost)
			router.Handle("/materials/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(matHandler.Edit)), // ✅ nuevo PUT
			).Methods(http.MethodPut)
			router.Handle("/materials/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(matHandler.Delete)),
			).Methods(http.MethodDelete)
		}

		// ======== AUDITS ========
		if s.AuditRepository != nil {
			auditHandler := handlers.NewAuditHandler(
				s.AuditRepository,
				s.HandleError,
				s.logger.Info,
			)
			router.HandleFunc("/audits", auditHandler.GetAll).Methods(http.MethodGet)
			router.HandleFunc("/audits/{id}", auditHandler.GetByID).Methods(http.MethodGet)
			router.Handle("/audits",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(auditHandler.Create)),
			).Methods(http.MethodPost)
			router.Handle("/audits/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(auditHandler.Edit)),
			).Methods(http.MethodPut)
			router.Handle("/audits/{id}",
				s.AuthMiddleware("supervisor")(http.HandlerFunc(auditHandler.Delete)),
			).Methods(http.MethodDelete)
		}

	}

	return router
}

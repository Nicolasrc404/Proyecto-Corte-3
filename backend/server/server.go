package server

import (
	"backend-avanzada/config"
	"backend-avanzada/logger"
	"backend-avanzada/models"
	"backend-avanzada/repository"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Server representa el servidor principal de la aplicación.
type Server struct {
	DB                      *gorm.DB
	Config                  *config.Config
	Handler                 http.Handler
	UserRepository          repository.UserRepository
	AlchemistRepository     *repository.AlchemistRepository
	MissionRepository       *repository.MissionRepository       // ✅ CRUD Missions
	MaterialRepository      *repository.MaterialRepository      // CRUD Materials
	TransmutationRepository *repository.TransmutationRepository // CRUD Transmutations
	AuditRepository         *repository.AuditRepository         // CRUD Audits
	jwtSecret               string
	logger                  *logger.Logger
	taskQueue               *TaskQueue
}

// NewServer inicializa la instancia del servidor.
func NewServer() *Server {
	s := &Server{
		logger: logger.NewLogger(),
	}
	var cfg config.Config
	configFile, err := os.ReadFile("config/config.json")
	if err != nil {
		s.logger.Fatal(err)
	}
	if err := json.Unmarshal(configFile, &cfg); err != nil {
		s.logger.Fatal(err)
	}
	s.Config = &cfg

	// Cargar secreto JWT desde .env
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		s.logger.Fatal(fmt.Errorf("JWT_SECRET is not set in environment"))
	}
	s.jwtSecret = secret

	return s
}

// StartServer arranca el servidor HTTP.
func (s *Server) StartServer() {
	fmt.Println("Inicializando base de datos...")
	s.initDB()
	if err := s.initAsyncInfrastructure(); err != nil {
		s.logger.Fatal(err)
	}

	fmt.Println("Configurando CORS...")
	corsObj := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
	)

	fmt.Println("Inicializando mux...")
	srv := &http.Server{
		Addr:    s.Config.Address,
		Handler: corsObj(s.router()),
	}
	fmt.Println("🚀 Escuchando en el puerto", s.Config.Address)
	if err := srv.ListenAndServe(); err != nil {
		s.logger.Fatal(err)
	}
}

// initDB configura la conexión a la base de datos y aplica migraciones.
func (s *Server) initDB() {
	switch s.Config.Database {
	case "sqlite":
		db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			s.logger.Fatal(err)
		}
		s.DB = db
	case "postgres":
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("POSTGRES_HOST"),
			os.Getenv("POSTGRES_USER"),
			os.Getenv("POSTGRES_PASSWORD"),
			os.Getenv("POSTGRES_DB"),
		)
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			s.logger.Fatal(err)
		}
		s.DB = db
	}

	fmt.Println("Aplicando migraciones...")

	// 🔹 Migraciones (sin borrar datos previos)
	err := s.DB.AutoMigrate(
		&models.User{},
		&models.Alchemist{},
		&models.Mission{}, // ✅ Importante para CRUD Missions
		&models.Material{},
		&models.Transmutation{},
		&models.Audit{},
	)
	if err != nil {
		s.logger.Fatal(err)
	}

	// 🔹 Inicializar repositorios
	s.UserRepository = repository.NewUserRepository(s.DB)
	s.AlchemistRepository = repository.NewAlchemistRepository(s.DB)
	s.MissionRepository = repository.NewMissionRepository(s.DB) // ✅
	s.MaterialRepository = repository.NewMaterialRepository(s.DB)
	s.TransmutationRepository = repository.NewTransmutationRepository(s.DB)
	s.AuditRepository = repository.NewAuditRepository(s.DB)
}
func (s *Server) initAsyncInfrastructure() error {
	redisAddr := s.Config.RedisAddress
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}
	s.taskQueue = NewTaskQueue(redisAddr, s.logger)
	s.taskQueue.WithRepositories(
		s.TransmutationRepository,
		s.AuditRepository,
		s.MissionRepository,
		s.MaterialRepository,
	)

	verificationInterval := time.Duration(s.Config.VerificationIntervalMinutes) * time.Minute
	pendingHours := time.Duration(s.Config.PendingTransmutationHours) * time.Hour
	lowStock := s.Config.MaterialLowStockThreshold

	s.taskQueue.ConfigureThresholds(verificationInterval, pendingHours, lowStock)
	if err := s.taskQueue.Start(); err != nil {
		return err
	}
	s.taskQueue.ScheduleDailyVerification()
	return nil
}

// GetJWTSecret devuelve la clave secreta usada para firmar los tokens JWT.
func (s *Server) GetJWTSecret() string {
	return s.jwtSecret
}

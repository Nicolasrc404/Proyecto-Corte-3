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

	"github.com/gorilla/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Server representa el servidor principal de la aplicación.
type Server struct {
	DB                  *gorm.DB
	Config              *config.Config
	Handler             http.Handler
	PeopleRepository    repository.Repository[models.Person]
	KillRepository      repository.Repository[models.Kill]
	UserRepository      repository.UserRepository
	AlchemistRepository *repository.AlchemistRepository // nuevo
	jwtSecret           string
	logger              *logger.Logger
	taskQueue           *TaskQueue
}

// NewServer inicializa la instancia del servidor.
func NewServer() *Server {
	s := &Server{
		logger:    logger.NewLogger(),
		taskQueue: NewTaskQueue(),
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

	// Migraciones (agrega nuevas entidades sin eliminar las anteriores)
	err := s.DB.AutoMigrate(
		&models.User{},
		&models.Person{},
		&models.Kill{},
		&models.Alchemist{}, // nuevo modelo
	)
	if err != nil {
		s.logger.Fatal(err)
	}

	// Inicializar repositorios
	s.KillRepository = repository.NewKillRepository(s.DB)
	s.PeopleRepository = repository.NewPeopleRepository(s.DB)
	s.UserRepository = repository.NewUserRepository(s.DB)

	// Solo si existe el modelo Alchemist
	s.AlchemistRepository = repository.NewAlchemistRepository(s.DB)
}

// GetJWTSecret devuelve la clave secreta usada para firmar los tokens JWT.
func (s *Server) GetJWTSecret() string {
	return s.jwtSecret
}

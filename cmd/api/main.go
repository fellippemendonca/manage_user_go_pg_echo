package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/healthz"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/messages"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/migrator"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/repositories"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/middlewares"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/routes"
)

type Todo struct {
	Name        string
	Description string
}

func main() {
	server := server.NewServer()

	// Init Zap Logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("logger initialization failed", err)
	}
	defer zapLogger.Sync()

	server.Logger = zapLogger
	server.Logger.Info("logger initialized")

	// Init .env variables
	err = godotenv.Load(".env")
	if err != nil {
		server.Logger.Warn("Error loading .env file") // May warn running inside docker
	}

	// Database Connection
	db, err := sql.Open("postgres", os.Getenv("MANAGE_USER_GO_POSTGRES"))
	if err != nil {
		server.Logger.Fatal("database connection failed", zap.Error(err))
	}
	defer db.Close()
	server.Logger.Info("Database connected")

	// Instantiating a new UserRepository
	userRepo := repositories.NewUserRepo(db)

	// Connecting to RabbitMQ
	conn, err := amqp.Dial(os.Getenv("MANAGE_USER_GO_RABBITMQ"))
	if err != nil {
		server.Logger.Fatal("rabbitmq connection failed", zap.Error(err))
	}
	defer conn.Close()

	// Opening Channel to RabbitMQ
	ch, err := conn.Channel()
	if err != nil {
		server.Logger.Fatal("rabbitmq open channel failed", zap.Error(err))
	}
	defer ch.Close()

	server.Logger.Info("Messaging service connected")

	// Assign Chained tester with active connections to Server
	server.ConnectionTester = &healthz.ChainedTester{
		Testers: []healthz.ConnectionTester{
			&healthz.DBTester{DB: db},
			&healthz.AmqpTester{Conn: conn},
		},
	}

	// Instantiating a new UserEvents wrapping UserRepository
	wrappedRepo := messages.NewUserEvents(server.Logger, ch, userRepo)

	// Assigning wrapped-UserRepository to Server
	server.UserRepository = wrappedRepo

	// Initial migrations if not yet exists
	if err := migrator.MigrateDB(os.Getenv("MANAGE_USER_GO_POSTGRES")); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			server.Logger.Info("migrations unchanged", zap.Error(err))
		} else {
			server.Logger.Fatal("migrations failed", zap.Error(err))
		}
	} else {
		server.Logger.Info("migrations done")
	}

	// Instantiating Echo
	e := echo.New()
	e.AcquireContext()

	// Logger middleware
	e.Use(middlewares.Logger(server.Logger))

	// Recover middleware
	e.Use(middlewares.Recover())

	// creating /api path group
	api := e.Group("/api")

	// Cors Middleware
	api.Use(middlewares.Cors())

	// Load Routes
	routes.LoadRoutes(api, server)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

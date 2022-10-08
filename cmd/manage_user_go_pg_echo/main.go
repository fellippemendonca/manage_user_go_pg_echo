package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/messages"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/repositories"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/healthz"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/middlewares"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/routes"
)

type Todo struct {
	Name        string
	Description string
}

func main() {

	server := server.NewServer()

	// Init .env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Init Zap Logger
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Fatal("Logger initialization failed", err)
	}
	defer zapLogger.Sync()

	server.Logger = zapLogger
	server.Logger.Info("Logger initialized")

	// Database Connection
	db, err := sql.Open("postgres", os.Getenv("MANAGE_USER_GO_POSTGRES"))
	if err != nil {
		server.Logger.Panic("Database connection failed", zap.Error(err))
	}
	defer db.Close()
	server.Logger.Info("Database connected")
	userRepo := repositories.NewUserRepo(db)

	conn, err := amqp.Dial(os.Getenv("MANAGE_USER_GO_RABBITMQ"))
	if err != nil {
		server.Logger.Panic("rabbitmq connection failed", zap.Error(err))
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		server.Logger.Panic("rabbitmq open channel failed", zap.Error(err))
	}
	defer ch.Close()

	server.Logger.Info("Messaging service connected")
	userEvents := messages.NewUserEvents(server.Logger, ch, userRepo)

	// Server
	server.ConnectionTester = &healthz.ChainedTester{
		Testers: []healthz.ConnectionTester{
			&healthz.DBTester{DB: db},
			&healthz.AmqpTester{Conn: conn},
		},
	}
	server.UserRepository = userEvents

	// Echo
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

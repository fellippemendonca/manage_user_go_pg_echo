package main

import (
	"database/sql"
	"log"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/messages"
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
		log.Fatal("Logger initialization failed", err)
	}
	defer zapLogger.Sync()

	server.Logger = zapLogger
	server.Logger.Info("Logger initialized")

	// Database Connection
	db, err := sql.Open("postgres", "postgresql://192.168.1.4/db?user=username&password=password&sslmode=disable")
	if err != nil {
		server.Logger.Panic("Database connection failed", zap.Error(err))
	}
	defer db.Close()

	userRepository := repositories.NewUserRepo(db)
	userEvents := messages.NewUserEvents(userRepository)
	server.UserRepository = userEvents
	server.Logger.Info("Database connected")

	// Echo
	e := echo.New()
	e.AcquireContext()

	// Logger middleware
	e.Use(middlewares.Logger(server.Logger))

	// Recover middleware
	e.Use(middlewares.Recover())

	// under /api path
	api := e.Group("/api")

	// Cors Middleware
	api.Use(middlewares.Cors())

	// Routes
	routes.LoadRoutes(api, server)

	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

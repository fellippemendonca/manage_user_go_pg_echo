package server

import (
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/healthz"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
)

// Server Struct is responsible to store the dependencies that will be used in the controllers
type Server struct {
	UserRepository   models.UserRepository
	Logger           *zap.Logger
	ConnectionTester healthz.ConnectionTester
}

func NewServer() *Server {
	return &Server{}
}

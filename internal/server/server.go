package server

import (
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"

	"go.uber.org/zap"
)

type Server struct {
	UserRepository models.UserRepository
	Logger         *zap.Logger
}

func NewServer() *Server {
	return &Server{}
}

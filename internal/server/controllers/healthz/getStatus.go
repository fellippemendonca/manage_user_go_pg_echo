package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
)

func GetStatus(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {

		err := s.UserRepository.TestConnection(c.Request().Context())
		if err != nil {
			s.Logger.Error("DB Connection check failed.", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}

		s.Logger.Info("DB Connection check successful.")
		return c.NoContent(http.StatusNoContent)
	}
}

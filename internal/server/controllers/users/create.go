package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
)

// Create User Controller is responsible for the User Creation
func Create(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := new(models.User)
		if err := c.Bind(u); err != nil {
			s.Logger.Error("failed to parse user body", zap.Error(err))
			return c.NoContent(http.StatusBadRequest)
		}

		user, err := s.UserRepository.CreateUser(c.Request().Context(), u)
		if err != nil {
			s.Logger.Error("failed to persist user", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusCreated, user)
	}
}

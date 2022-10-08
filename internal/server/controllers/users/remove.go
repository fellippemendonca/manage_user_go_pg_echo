package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
)

// e.DELETE("/users/:id", remove)
func Remove(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")

		parsedID, err := uuid.Parse(id)
		if err != nil {
			s.Logger.Error("failed to parse user id", zap.Error(err))
			return c.NoContent(http.StatusBadRequest)
		}

		_, err = s.UserRepository.RemoveUser(c.Request().Context(), parsedID)
		if err != nil {
			s.Logger.Error("failed to remove user", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusAccepted)
	}
}

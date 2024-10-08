package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
)

// Remove User Controller is responsible for Deleting the user from the database using its ID.
func Remove(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		parsedID, err := uuid.Parse(c.Param("id"))
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

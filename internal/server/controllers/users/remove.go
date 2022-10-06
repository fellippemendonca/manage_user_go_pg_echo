package controllers

import (
	"net/http"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// e.DELETE("/users/:id", remove)
func Remove(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {

		id := c.Param("id")

		_, err := s.UserRepository.RemoveUser(c.Request().Context(), uuid.Must(uuid.Parse(id)))
		if err != nil {
			s.Logger.Error("Failed to find users")
			return err
		}

		return c.NoContent(http.StatusAccepted)
	}
}

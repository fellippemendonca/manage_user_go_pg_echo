package users

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
)

// e.PUT("/users", putUser)
func Update(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		parsedID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			s.Logger.Error("failed to parse id", zap.Error(err))
			return c.NoContent(http.StatusBadRequest)
		}

		u := new(models.User)
		if err := c.Bind(u); err != nil {
			s.Logger.Error("failed to parse user body", zap.Error(err))
			return c.NoContent(http.StatusBadRequest)
		}

		if parsedID != u.ID {
			s.Logger.Error("param id does not match body id")
			return c.NoContent(http.StatusBadRequest)
		}

		user, err := s.UserRepository.UpdateUser(c.Request().Context(), u)
		if err != nil {
			s.Logger.Error("failed to update user", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, user)
	}
}

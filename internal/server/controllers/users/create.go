package controllers

import (
	"net/http"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"

	"github.com/labstack/echo/v4"
)

// e.POST("/users", postUser)
func Create(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {

		u := new(models.User)
		if err := c.Bind(u); err != nil {
			s.Logger.Error("Failed to parse user body")
			return err
		}

		res, err := s.UserRepository.CreateUser(c.Request().Context(), u)
		if err != nil {
			s.Logger.Error("Failed to persist user")
			return err
		}

		return c.JSON(http.StatusCreated, res)
	}
}

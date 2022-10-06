package controllers

import (
	"net/http"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"

	"github.com/labstack/echo/v4"
)

// e.PUT("/users", putUser)
func Update(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		u := new(models.User)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	}
}

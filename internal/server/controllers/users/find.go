package controllers

import (
	"net/http"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"

	"github.com/labstack/echo/v4"
)

// e.GET("/users", list)
func Find(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {

		user := models.User{}
		values := c.Request().URL.Query()

		for k, v := range values {
			switch k {
			case "country":
				user.Country = v[0]
			case "first_name":
				user.First_name = v[0]
			case "last_name":
				user.Last_name = v[0]
			case "email":
				user.Email = v[0]
			case "nickname":
				user.Nickname = v[0]
			}
		}

		result, err := s.UserRepository.FindUsers(c.Request().Context(), &user)
		if err != nil {
			s.Logger.Error("Failed to find users")
			return err
		}

		return c.JSON(http.StatusOK, result)
	}
}

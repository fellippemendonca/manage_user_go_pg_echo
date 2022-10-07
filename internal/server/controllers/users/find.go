package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	"go.uber.org/zap"

	"github.com/labstack/echo/v4"
)

// e.GET("/users", list)
func Find(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		values := c.Request().URL.Query()

		var limit int
		var err error
		if values.Has("limit") {
			limitStr := values.Get("limit")
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				s.Logger.Error("Failed to parse page limit", zap.Error(err))
				return c.NoContent(http.StatusBadRequest)
			}
		}

		user := models.User{}
		user.Country = values.Get("country")
		user.First_name = values.Get("first_name")
		user.Last_name = values.Get("last_name")
		user.Email = values.Get("email")
		user.Nickname = values.Get("nickname")
		pageToken := values.Get("page_token")

		users, pageToken, err := s.UserRepository.FindUsers(c.Request().Context(), &user, pageToken, limit)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			s.Logger.Error("FindUsers failed", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, models.UsersResponse{Users: users, PageToken: pageToken})
	}
}

package users

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
)

// Find Users Controller is able to retrieve a paginated list of users based on the query used and page-limit.
func Find(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		values := c.Request().URL.Query()

		var limit int
		var err error
		// Check and convert page-limit from querystring
		if values.Has("limit") {
			limitStr := values.Get("limit")
			limit, err = strconv.Atoi(limitStr)
			if err != nil {
				s.Logger.Error("Failed to parse page limit", zap.Error(err))
				return c.NoContent(http.StatusBadRequest)
			}
		}

		// Fills a User object with possible values provided in the querystring
		user := &models.User{}
		user.Country = values.Get("country")
		user.FirstName = values.Get("first_name")
		user.LastName = values.Get("last_name")
		user.Email = values.Get("email")
		user.Nickname = values.Get("nickname")
		pageToken := values.Get("page_token")

		usersResponse, err := s.UserRepository.FindUsers(c.Request().Context(), user, pageToken, limit)
		if err != nil {
			s.Logger.Error("FindUsers failed", zap.Error(err))
			if errors.Is(err, sql.ErrNoRows) {
				return c.NoContent(http.StatusNotFound)
			}
			return c.NoContent(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, usersResponse)
	}
}

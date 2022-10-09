package healthz

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
)

// GetStatus is the controller responsible to trigger all dependencies connection tests.
func GetStatus(s *server.Server) func(c echo.Context) error {
	return func(c echo.Context) error {
		err := s.ConnectionTester.TestConnection(c.Request().Context())
		if err != nil {
			s.Logger.Error("health check failed", zap.Error(err))
			return c.NoContent(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

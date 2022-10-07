package routes

import (
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	healthz "github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/healthz"
	users "github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/users"

	"github.com/labstack/echo/v4"
)

func LoadRoutes(g *echo.Group, s *server.Server) {
	g.GET("/healthz", healthz.GetStatus(s))

	g.GET("/users", users.Find(s))
	g.POST("/users", users.Create(s))
	g.PUT("/users", users.Update(s))
	g.DELETE("/users/:id", users.Remove(s))
}

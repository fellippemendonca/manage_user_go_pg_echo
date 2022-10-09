package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/healthz"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/users"
)

func LoadRoutes(g *echo.Group, s *server.Server) {
	g.GET("/healthz", healthz.GetStatus(s))

	g.GET("/users", users.Find(s))
	g.POST("/users", users.Create(s))
	g.PUT("/users/:id", users.Update(s))
	g.DELETE("/users/:id", users.Remove(s))
}

package routes

import (
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	healthz "github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/healthz"
	users "github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/users"

	"github.com/labstack/echo/v4"
)

func LoadRoutes(g *echo.Group, s *server.Server) {

	g.GET("/healthz", healthz.GetStatus(s))

	v1 := g.Group("/v1")
	v1.GET("/users", users.Find(s))
	v1.POST("/users", users.Create(s))
	v1.PUT("/users/:id", users.Update(s))
	v1.DELETE("/users/:id", users.Remove(s))
}

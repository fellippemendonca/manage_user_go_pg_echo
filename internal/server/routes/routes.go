package routes

import (
	"github.com/labstack/echo/v4"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/healthz"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/users"
)

// LoadRoutes is responsible to assign the paths to the methods and also assign the Server to the Controllers
func LoadRoutes(g *echo.Group, s *server.Server) {
	g.GET("/healthz", healthz.GetStatus(s))

	g.GET("/users", users.Find(s)) // There is no get by ID because it wasn't in the requirements.
	g.POST("/users", users.Create(s))
	g.PUT("/users/:id", users.Update(s)) // Should not be used as PATCH! All User fields shold be provided otherwise will be blanked.
	g.DELETE("/users/:id", users.Remove(s))
}

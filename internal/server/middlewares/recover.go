package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Recover middleware is just using the Echo existing one.
func Recover() echo.MiddlewareFunc {
	return middleware.Recover()
}

package routes

import (
	"github.com/genpsp/go-app/services/src/handler"
	"github.com/genpsp/go-app/services/src/middlewares"
	"github.com/labstack/echo/v4"
)

func Init(handler handler.Handler, m middlewares.Middleware, e *echo.Echo) {
	admin := e.Group("/app")

	item := admin.Group("/item")
	item.GET("", handler.Enum.GetEnums, m.Auth.RequireJWTAuthorizationHeader())

}

package main

import (
	"github.com/genpsp/go-app/pkg/channel"
	"github.com/genpsp/go-app/pkg/configs"
	"github.com/genpsp/go-app/pkg/database"
	"github.com/genpsp/go-app/pkg/firebase"
	"github.com/genpsp/go-app/pkg/logger"
	"github.com/genpsp/go-app/pkg/server"
	"github.com/genpsp/go-app/services/src/handler"
	"github.com/genpsp/go-app/services/src/middlewares"
	"github.com/genpsp/go-app/services/src/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	configs.LoadConfig()
	cfg := configs.GetConfig()
	logger.LoadLogger(cfg.System.Env, cfg.Logger.LogLevel, cfg.Logger.LogEncoding)

	db := database.Open(cfg.MySQL)
	defer db.Close()

	httpServer := server.NewHttpServer()
	authClient := firebase.NewFirebaseAppAdmin()
	handler := handler.NewHandler(db.Master, authClient)
	middleware := middlewares.NewMiddleware(authClient)

	httpServer.Handler = func(e *echo.Echo) {
		routes.Init(handler, middleware, e)
	}

	httpServer.Start()

	signal := <-channel.Quit()
	defer httpServer.Stop(signal)
}

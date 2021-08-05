package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4/middleware"
	"os"
	"time"

	appErr "github.com/genpsp/go-app/pkg/server/error"

	"github.com/genpsp/go-app/pkg/channel"
	"github.com/genpsp/go-app/pkg/configs"
	"github.com/genpsp/go-app/pkg/logger"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
)

type HttpServer struct {
	echo    *echo.Echo
	Handler func(server *echo.Echo)
	Addr    string
	Timeout time.Duration
}

type Context struct {
	echo.Context
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewHttpServer() HttpServer {
	config := configs.GetConfig()
	e := echo.New()
	e.Server.Addr = fmt.Sprintf(":%s", config.System.HttpAddr)
	e.HideBanner = true
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Use(middleware.CORS())
	e.HTTPErrorHandler = appErr.JSONErrorHandler
	loc, _ := time.LoadLocation(config.System.TimeZone)
	logger.Logging.Info(fmt.Sprintf("current timezone: %s", loc))

	return HttpServer{
		echo:    e,
		Addr:    fmt.Sprintf(":%s", config.System.HttpAddr),
		Timeout: config.System.HttpContextTimeoutSec,
	}
}

func (srv *HttpServer) Start() {
	srv.Handler(srv.echo)
	go func() {
		logger.Logging.Info(fmt.Sprintf("start http server addr: %s", srv.Addr))
		if err := srv.echo.StartServer(srv.echo.Server); err != nil {
			logger.Logging.Error(fmt.Sprintf("http serve error: %s", err.Error()))
		}
	}()
}

func (srv *HttpServer) Stop(signal os.Signal) {
	ctx, cancel := context.WithTimeout(context.Background(), srv.Timeout*time.Second)
	defer cancel()
	if err := srv.echo.Shutdown(ctx); err != nil {
		srv.echo.Logger.Fatal(err)
	}
	logger.Logging.Info(fmt.Sprintf("stopping http server... ExitCode: %d, Signal: %s", channel.GetExitCode(signal), signal.String()))
}

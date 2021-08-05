package middlewares

import (
	"fmt"
	"github.com/genpsp/go-app/domain/enum"
	"github.com/genpsp/go-app/pkg/logger"
	"github.com/genpsp/go-app/services/src/services"

	appErr "github.com/genpsp/go-app/pkg/server/error"

	"github.com/genpsp/go-app/pkg/server/jwt"
	"github.com/labstack/echo/v4"
)

type (
	Auth interface {
		RequireJWTAuthorizationHeader() echo.MiddlewareFunc
	}

	authImpl struct {
		as services.AuthService
	}
)

func NewAuth(s services.AuthService) Auth {
	return &authImpl{
		as: s,
	}
}

type JWTContext struct {
	echo.Context
	Token *jwt.Token
}

func (s *authImpl) RequireJWTAuthorizationHeader() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			firebaseJWT := c.Request().Header.Get("Authorization")
			token, err := s.as.Authorize(firebaseJWT)
			if err != nil {
				return appErr.BindAppErrorWithServiceError(err)
			}

			jc := &JWTContext{c, &jwt.Token{UID: token.UID}}

			c.Set("token", jc.Token)
			return next(jc)
		}
	}
}

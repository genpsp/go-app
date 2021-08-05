package middlewares

import (
	"github.com/genpsp/go-app/pkg/firebase"
	"github.com/genpsp/go-app/services/src/services"
)

type (
	Middleware struct {
		Auth Auth
	}
)

func NewMiddleware(authClient firebase.AuthAdmin) Middleware {
	// service
	authService := services.NewAuthService(&authClient)

	return Middleware{
		Auth: NewAuth(authService),
	}
}

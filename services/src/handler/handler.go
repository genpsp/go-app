package handler

import (
	repositories "github.com/genpsp/go-app/domain/repository"
	"github.com/genpsp/go-app/pkg/configs/cloudfunctions"
	"github.com/genpsp/go-app/pkg/configs/gcs"
	"github.com/genpsp/go-app/pkg/firebase"
	"github.com/genpsp/go-app/services/src/services"
	"gorm.io/gorm"
)

type (
	Handler struct {
		Item Item
	}
)

func NewHandler(m *gorm.DB, f firebase.AuthAdmin) Handler {
	// repository
	itemRepo := repositories.NewItemRepository()

	// service
	itemService := services.NewItemService(itemRepo, m, f)

	return Handler{
		Item: NewItem(itemService, f),
	}
}

package services

import (
	"fmt"
	entities "github.com/genpsp/go-app/domain/entities"
	repositories "github.com/genpsp/go-app/domain/repository"
	"github.com/genpsp/go-app/pkg/firebase"
	"github.com/genpsp/go-app/pkg/logger"
	appErr "github.com/genpsp/go-app/pkg/server/error"
	"github.com/genpsp/go-app/services/src/handler/request"
	"gorm.io/gorm"
)

type (
	ItemService interface {
		FindAll() (items *[]entities.Item, err error)
		Find(gar *request.GetItemRequest) (items *[]entities.Item, err error)
		FindByID(itemID int) (item *entities.Item, err error)
		Create(itemEntity *entities.Item, password string) (err error)
		Update(itemID int, itemEntity *entities.Item) (err error)
		Delete(itemID int) (err error)
	}

	itemServiceImpl struct {
		aur    repositories.ItemRepository
		master *gorm.DB
		auth   firebase.AuthAdmin
	}
)

func NewItemService(
	itemRepo repositories.ItemRepository,
	m *gorm.DB, auth firebase.AuthAdmin) ItemService {

	return &itemServiceImpl{
		aur:    itemRepo,
		master: m,
		auth:   auth,
	}
}

func (s *itemServiceImpl) FindAll() (items *[]entities.Item, err error) {
	err = s.master.Transaction(func(tx *gorm.DB) error {
		items, err = s.aur.FindAll(tx)
		if err != nil {
			logger.Logging.Error(fmt.Sprintf("occurred error when Item with FindAll call ItemRepository: %s", err.Error()))
			return appErr.BindServiceErrorWithDBErrorCaseRecordNotFoundIsNil(err)
		}
		return nil
	})
	return
}

func (s *itemServiceImpl) Find(gar *request.GetItemRequest) (items *[]entities.Item, err error) {
	err = s.master.Transaction(func(tx *gorm.DB) error {
		items, err = s.aur.Find(tx, gar.EmailAddress, gar.Name)
		if err != nil {
			logger.Logging.Error(fmt.Sprintf("occurred error when Item with FindOne call ItemRepository: %s", err.Error()))
			return appErr.BindServiceErrorWithDBErrorCaseRecordNotFoundIsNil(err)
		}
		return nil
	})
	return
}

func (s *itemServiceImpl) FindByID(itemID int) (itemEntity *entities.Item, err error) {
	err = s.master.Transaction(func(tx *gorm.DB) error {
		itemEntity, err = s.aur.FindByID(tx, itemID)
		if err != nil {
			logger.Logging.Error(fmt.Sprintf("occurred error when Item with FindByID call ItemRepository: %s", err.Error()))
			return appErr.BindServiceErrorWithDBError(err)
		}
		return nil
	})
	return
}

func (s *itemServiceImpl) Create(itemEntity *entities.Item, password string) (err error) {
	err = s.master.Transaction(func(tx *gorm.DB) error {
		result, createUserErr := s.auth.CreateUser(itemEntity, password)
		if createUserErr != nil || result == nil {
			return appErr.BindServiceErrorWithDBError(createUserErr)
		}
		itemEntity.ExternalUserID = result.UID

		createErr := s.aur.Create(tx, itemEntity)
		if createErr != nil {
			deleteUserErr := s.auth.DeleteUser(result.UID)
			return appErr.BindServiceErrorWithDBError(deleteUserErr)
		}

		claims := map[string]interface{}{"role": itemEntity.Role}
		setCustomClaimsErr := s.auth.SetCustomClaims(itemEntity.ExternalUserID, claims)
		if setCustomClaimsErr != nil {
			deleteUserErr := s.auth.DeleteUser(result.UID)
			return appErr.BindServiceErrorWithFirebaseError(deleteUserErr)
		}
		return nil
	})
	return
}

func (s *itemServiceImpl) Update(itemID int, itemEntity *entities.Item) (err error) {
	err = s.master.Transaction(func(tx *gorm.DB) error {
		err := s.aur.Update(tx, itemID, itemEntity)
		if err != nil {
			logger.Logging.Error(fmt.Sprintf("occurred error when Item with Update call ItemRepository: %s", err.Error()))
			return appErr.BindServiceErrorWithDBError(err)
		}
		return nil
	})
	return
}

func (s *itemServiceImpl) Delete(itemID int) (err error) {
	err = s.master.Transaction(func(tx *gorm.DB) error {
		item, err := s.aur.FindByID(tx, itemID)
		if err != nil {
			logger.Logging.Error(fmt.Sprintf("occurred error when Item with Delete call ItemRepository: %s", err.Error()))
			return appErr.BindServiceErrorWithDBError(err)
		}

		if err = s.auth.DeleteUser(item.ExternalUserID); err != nil {
			logger.Logging.Error(fmt.Sprintf("occurred error when Item with Delete call Firebase deleteUser: %s", err.Error()))
			return appErr.BindServiceErrorWithFirebaseError(err)
		}

		err = s.aur.Delete(tx, itemID)
		if err != nil {
			logger.Logging.Error(fmt.Sprintf("occurred error when Item with Delete call ItemRepository: %s", err.Error()))
			return appErr.BindServiceErrorWithDBError(err)
		}
		return nil
	})
	return
}

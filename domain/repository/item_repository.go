package repositories

import (
	"errors"
	"fmt"

	entities "github.com/genpsp/go-app/domain/entities"
	"github.com/genpsp/go-app/pkg/logger"
	appErr "github.com/genpsp/go-app/pkg/server/error"
	"gorm.io/gorm"
)

type (
	ItemRepository interface {
		FindAll(db *gorm.DB) (items *[]entities.Item, err error)
		FindByID(db *gorm.DB, itemID int) (itemEntity *entities.Item, err error)
		Create(db *gorm.DB, itemEntity *entities.Item) (err error)
		Update(db *gorm.DB, itemID int, itemEntity *entities.Item) (err error)
		Delete(db *gorm.DB, itemID int) (err error)
	}
	ItemRepositoryImpl struct{}
)

func NewItemRepository() ItemRepository {
	return &ItemRepositoryImpl{}
}

func (r *ItemRepositoryImpl) FindAll(db *gorm.DB) (items *[]entities.Item, err error) {
	err = db.Model(&entities.Item{}).
		Find(&items).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logging.Info(fmt.Sprintf("Item. record not found."))
		return nil, nil
	}

	if err != nil {
		logger.Logging.Error(fmt.Sprintf("Item FindAll error: %s", err.Error()))
		err = appErr.DBClientError
		return
	}

	return
}

func (r *ItemRepositoryImpl) FindByID(db *gorm.DB, itemID int) (itemEntity *entities.Item, err error) {
	err = db.Model(&entities.Item{}).
		Where("id = ?", itemID).
		First(&itemEntity).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Logging.Info(fmt.Sprintf("Item. record not found."))
		return nil, nil
	}

	if err != nil {
		logger.Logging.Error(fmt.Sprintf("Item FindByID error: %s", err.Error()))
		err = appErr.DBClientError
		return
	}

	return
}

func (r *ItemRepositoryImpl) Create(db *gorm.DB, itemEntity *entities.Item) (err error) {
	err = db.Create(&itemEntity).Error

	if err != nil {
		logger.Logging.Error(fmt.Sprintf("Item Create error: %s", err.Error()))
		err = appErr.DBClientError
		return
	}

	return
}

func (r *ItemRepositoryImpl) Update(db *gorm.DB, itemID int, itemEntity *entities.Item) (err error) {
	err = db.Model(&itemEntity).
		Where("id = ?", itemID).
		Updates(entities.Item{
			Name: itemEntity.Name,
		}).
		Error

	if err != nil {
		logger.Logging.Error(fmt.Sprintf("Item Update error: %s", err.Error()))
		err = appErr.DBClientError
		return
	}

	return
}

func (r *ItemRepositoryImpl) Delete(db *gorm.DB, itemID int) (err error) {
	itemEntity := entities.Item{}
	err = db.Model(&itemEntity).Where("id = ?", itemID).Delete(&itemEntity).Error
	if err != nil {
		logger.Logging.Error(fmt.Sprintf("Item Update error: %s", err.Error()))
		err = appErr.DBClientError
		return
	}
	return
}

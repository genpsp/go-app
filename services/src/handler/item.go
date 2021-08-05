package handler

import (
	"fmt"
	"net/http"
	"strconv"

	admin_response "github.com/genpsp/go-app/services/src/handler/response"

	entities "github.com/genpsp/go-app/domain/entities"

	"github.com/genpsp/go-app/pkg/firebase"

	"github.com/genpsp/go-app/pkg/logger"
	appErr "github.com/genpsp/go-app/pkg/server/error"
	"github.com/genpsp/go-app/pkg/utils"
	"github.com/genpsp/go-app/services/src/handler/request"
	"github.com/genpsp/go-app/services/src/services"
	"github.com/labstack/echo/v4"
)

type (
	Item interface {
		Find(c echo.Context) (err error)
		FindByID(c echo.Context) (err error)
		Create(c echo.Context) (err error)
		Update(c echo.Context) (err error)
		Delete(c echo.Context) (err error)
	}
	itemImpl struct {
		aus  services.ItemService
		auth firebase.AuthAdmin
	}
)

func NewItem(s services.ItemService, f firebase.AuthAdmin) Item {
	return &itemImpl{
		aus:  s,
		auth: f,
	}
}

func (s *itemImpl) Find(c echo.Context) (err error) {
	gar := new(request.GetItemRequest)
	if _, err := utils.RequestValidate(c, gar); err != "" {
		logger.Logging.Error(fmt.Sprintf("parse in GetItem erros: %s,  body: %s", err, utils.ToJson(gar)))
		return appErr.AppStatusBadRequestError400
	}
	var result *[]entities.Item
	result, err = s.aus.FindAll()
	if err != nil {
		return appErr.BindAppErrorWithServiceError(err)
	}
	if result == nil {
		c.JSON(http.StatusNoContent, nil)
		return nil
	}
	items := admin_response.ConvertItemsResponse(result)
	c.JSON(http.StatusOK, items)
	return nil
}

func (s *itemImpl) FindByID(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("itemId"))
	result, err := s.aus.FindByID(id)
	if err != nil {
		return appErr.BindAppErrorWithServiceError(err)
	}
	if result == nil {
		c.JSON(http.StatusNoContent, nil)
		return nil
	}
	itemResponse := admin_response.ConvertItemResponse(*result)
	c.JSON(http.StatusOK, itemResponse)
	return nil
}

func (s *itemImpl) Create(c echo.Context) (err error) {
	car := new(request.CreateItemRequest)
	if _, err := utils.RequestValidate(c, car); err != "" {
		logger.Logging.Error(fmt.Sprintf("parse in CreateItemRequest erros: %s,  body: %s", err, utils.ToJson(car)))
		return appErr.AppStatusBadRequestError400
	}

	firebaseJwt := c.Request().Header.Get("Authorization")
	if firebaseJwt == "" {
		c.JSON(http.StatusBadRequest, appErr.AppStatusBadRequestError400)
		return
	}

	_, err = s.auth.VerifyIDToken(firebaseJwt)
	if err != nil {
		c.JSON(http.StatusUnauthorized, appErr.BindAppErrorWithServiceError(err))
		return
	}

	entity := &entities.Item{
		Name: car.Name,
	}
	password := utils.RandomString(8)
	if err = s.aus.Create(entity, password); err != nil {
		return appErr.AppStatusBadRequestError400
	}

	c.JSON(http.StatusCreated, nil)
	return
}

func (s *itemImpl) Update(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("itemId"))
	car := new(request.CreateItemRequest)
	if _, err := utils.RequestValidate(c, car); err != "" {
		logger.Logging.Error(fmt.Sprintf("parse in CreateItemRequest erros: %s,  body: %s", err, utils.ToJson(car)))
		return appErr.AppStatusBadRequestError400
	}

	entity := &entities.Item{
		Name: car.Name,
	}
	err = s.aus.Update(id, entity)
	if err != nil {
		return appErr.BindAppErrorWithServiceError(err)
	}
	c.JSON(http.StatusCreated, nil)
	return
}

func (s *itemImpl) Delete(c echo.Context) (err error) {
	id, _ := strconv.Atoi(c.Param("itemId"))
	err = s.aus.Delete(id)
	if err != nil {
		return appErr.BindAppErrorWithServiceError(err)
	}
	c.JSON(http.StatusNoContent, nil)
	return
}

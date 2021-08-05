package handler

import (
	"encoding/json"
	appErr "github.com/genpsp/go-app/pkg/server/error"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	admin_response "github.com/genpsp/go-app/services/src/handler/response"

	entities "github.com/genpsp/go-app/domain/entities"

	"github.com/genpsp/go-app/pkg/utils"
	"github.com/genpsp/go-app/services/src/handler/request"

	mock_services "github.com/genpsp/go-app/services/src/services/mock"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewItemHandler(t *testing.T) {
	Convey("ItemHandlerを初期化", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		as := mock_services.NewMockItemService(ctrl)
		ah := NewItem(as, nil)
		So(ah, ShouldNotBeNil)
	})
}

func Test_ItemHandler(t *testing.T) {
	Convey("ItemHandlerを初期化", t, func() {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		const itemID = 1
		const name = "テスト"
		const emailAddress = "test@gmail.com"
		const role = 0

		as := mock_services.NewMockItemService(ctrl)
		ah := NewItem(as, nil)
		So(ah, ShouldNotBeNil)

		Convey("FindAll", func() {
			e := echo.New()
			e.Validator = utils.NewAppValidator()
			req := httptest.NewRequest(http.MethodGet, "/admin_users?emailAddress=&name=", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockEntities := []entities.Item{
				{Name: name, EmailAddress: emailAddress, Role: role},
			}
			as.EXPECT().FindAll().Return(&mockEntities, nil)

			Convey("正常にレスポンスを変換できる", func() {
				response := admin_response.ConvertItemsResponse(&mockEntities)
				mockResponse := []*admin_response.ItemResponse{{
					Name: name, EmailAddress: emailAddress, Role: role,
				}}
				So(response, ShouldResemble, mockResponse)
				Convey("正常にレスポンスが返る", func() {
					err := ah.Find(c)
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusOK)
				})
			})
		})
		Convey("Find", func() {
			e := echo.New()
			e.Validator = utils.NewAppValidator()
			req := httptest.NewRequest(http.MethodGet, "/admin_users?emailAddress="+emailAddress+"&name="+name, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockEntities := []entities.Item{
				{Name: name, EmailAddress: emailAddress, Role: role},
			}
			mockRequest := request.GetItemRequest{
				Name: name, EmailAddress: emailAddress,
			}
			as.EXPECT().Find(&mockRequest).Return(&mockEntities, nil)

			Convey("正常にレスポンスを変換できる", func() {
				response := admin_response.ConvertItemsResponse(&mockEntities)
				mockResponse := []*admin_response.ItemResponse{{
					Name: name, EmailAddress: emailAddress, Role: role,
				}}
				So(response, ShouldResemble, mockResponse)
				Convey("正常にレスポンスが返る", func() {
					err := ah.Find(c)
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusOK)
				})
			})
		})
		Convey("FindByID", func() {
			e := echo.New()
			e.Validator = utils.NewAppValidator()
			req := httptest.NewRequest(http.MethodGet, "/admin_users/:itemId", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockEntity := entities.Item{
				Name: name, EmailAddress: emailAddress, Role: role,
			}
			as.EXPECT().FindByID(gomock.Any()).Return(&mockEntity, nil)

			Convey("正常にレスポンスを変換できる", func() {
				response := admin_response.ConvertItemResponse(mockEntity)
				mockResponse := admin_response.ItemResponse{
					Name: name, EmailAddress: emailAddress, Role: role,
				}
				So(response, ShouldResemble, &mockResponse)
				Convey("正常にレスポンスが返る", func() {
					err := ah.FindByID(c)
					So(err, ShouldBeNil)
					So(rec.Code, ShouldEqual, http.StatusOK)
				})
			})
		})
		Convey("Update", func() {
			e := echo.New()
			e.Validator = utils.NewAppValidator()

			body := request.CreateItemRequest{
				Name:         name,
				EmailAddress: emailAddress,
				Role:         role,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPut, "/admin_users", strings.NewReader(string(jsonBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			mockEntity := &entities.Item{
				Name: name, EmailAddress: emailAddress, Role: role,
			}

			as.EXPECT().Update(gomock.Any(), mockEntity).Return(nil)

			Convey("正常に更新できる", func() {
				err := ah.Update(c)
				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, http.StatusCreated)
			})
		})
		Convey("Delete", func() {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/admin_users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			as.EXPECT().Delete(gomock.Any()).Return(nil)

			Convey("正常に削除できる", func() {
				err := ah.Delete(c)
				So(err, ShouldBeNil)
				So(rec.Code, ShouldEqual, http.StatusNoContent)
			})
		})
		Convey("Findで正常に取得できなかった場合エラーを返す", func() {
			e := echo.New()
			e.Validator = utils.NewAppValidator()
			req := httptest.NewRequest(http.MethodGet, "/admin_users?emailAddress="+emailAddress+"&name="+name, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			as.EXPECT().Find(gomock.Any()).Return(nil, appErr.ServiceClientError)

			err := ah.Find(c)
			So(err, ShouldEqual, appErr.AppStatusInternalServerError500)
		})
		Convey("FindByIDで正常に取得できなかった場合エラーを返す", func() {
			e := echo.New()
			e.Validator = utils.NewAppValidator()
			req := httptest.NewRequest(http.MethodGet, "/admin_users/:itemId", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			as.EXPECT().FindByID(gomock.Any()).Return(nil, appErr.ServiceClientError)

			err := ah.FindByID(c)
			So(err, ShouldEqual, appErr.AppStatusInternalServerError500)
		})
		Convey("Updateで正常に更新できなかった場合エラーを返す", func() {
			e := echo.New()
			e.Validator = utils.NewAppValidator()

			body := request.CreateItemRequest{
				Name:         name,
				EmailAddress: emailAddress,
				Role:         role,
			}
			jsonBody, _ := json.Marshal(body)

			req := httptest.NewRequest(http.MethodPut, "/admin_users", strings.NewReader(string(jsonBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			as.EXPECT().Update(gomock.Any(), gomock.Any()).Return(appErr.ServiceClientError)

			err := ah.Update(c)
			So(err, ShouldEqual, appErr.AppStatusInternalServerError500)
		})
		Convey("Deleteで正常に削除できなかった場合エラーを返す", func() {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/admin_users", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			as.EXPECT().Delete(gomock.Any()).Return(appErr.ServiceClientError)

			err := ah.Delete(c)
			So(err, ShouldEqual, appErr.AppStatusInternalServerError500)
		})
	})
}

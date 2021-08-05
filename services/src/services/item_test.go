package services

import (
	"firebase.google.com/go/v4/auth"
	"github.com/genpsp/go-app/pkg/configs"
	"github.com/genpsp/go-app/pkg/logger"
	"github.com/genpsp/go-app/pkg/mock_pkgs"
	appErr "github.com/genpsp/go-app/pkg/server/error"
	"github.com/genpsp/go-app/pkg/utils"
	"github.com/genpsp/go-app/services/src/handler/request"
	"testing"

	entities "github.com/genpsp/go-app/domain/entities"
	"github.com/genpsp/go-app/domain/repository/mock_repositories"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewItemService(t *testing.T) {
	Convey("ItemServiceを初期化", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		db, _, _ := mock_repositories.GetDBMock()
		ar := mock_repositories.NewMockItemRepository(ctrl)
		fbAuth := mock_pkgs.NewMockAuthAdmin(ctrl)
		as := NewItemService(ar, db, fbAuth)
		So(as, ShouldNotBeNil)
	})
}

func Test_ItemService(t *testing.T) {
	Convey("ItemServiceを初期化", t, func() {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		configs.TestLoadConfig()
		cfg := configs.GetConfig()
		logger.LoadLogger(cfg.System.Env, cfg.Logger.LogLevel, cfg.Logger.LogEncoding)

		db, mock, _ := mock_repositories.GetDBMock()
		const itemID = 1
		const name = "テスト"
		const externalUserID = "1"
		const emailAddress = "test@gmail.com"
		const role = 0
		var password = utils.RandomString(8)

		ar := mock_repositories.NewMockItemRepository(ctrl)
		fbAuth := mock_pkgs.NewMockAuthAdmin(ctrl)
		as := NewItemService(ar, db, fbAuth)
		So(as, ShouldNotBeNil)

		Convey("FindAll", func() {
			mockEntities := []entities.Item{
				{Name: name, EmailAddress: emailAddress, Role: role},
			}
			mock.ExpectBegin()
			ar.EXPECT().FindAll(gomock.Any()).Return(&mockEntities, nil)
			mock.ExpectCommit()
			Convey("正常に取得できる", func() {
				result, err := as.FindAll()
				So(result, ShouldResemble, &mockEntities)
				So(err, ShouldBeNil)
			})
		})
		Convey("Find", func() {
			mockEntities := []entities.Item{{
				Name: name, EmailAddress: emailAddress, Role: role,
			}}
			mockRequest := request.GetItemRequest{
				Name: name, EmailAddress: emailAddress,
			}
			mock.ExpectBegin()
			ar.EXPECT().Find(gomock.Any(), emailAddress, name).Return(&mockEntities, nil)
			mock.ExpectCommit()
			Convey("正常に取得できる", func() {
				result, err := as.Find(&mockRequest)
				So(result, ShouldResemble, &mockEntities)
				So(err, ShouldBeNil)
			})
		})
		Convey("FindByID", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}
			mock.ExpectBegin()
			ar.EXPECT().FindByID(gomock.Any(), itemID).Return(mockEntity, nil)
			mock.ExpectCommit()
			Convey("正常に取得できる", func() {
				result, err := as.FindByID(itemID)
				So(result, ShouldResemble, mockEntity)
				So(err, ShouldBeNil)
			})
		})
		Convey("Create", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}
			mockResult := &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					DisplayName: name, Email: emailAddress, UID: externalUserID,
				},
				CustomClaims:  map[string]interface{}{"role": role},
				Disabled:      false,
				EmailVerified: false,
			}

			mock.ExpectBegin()
			ar.EXPECT().Create(gomock.Any(), mockEntity).Return(nil)
			fbAuth.EXPECT().CreateUser(mockEntity, password).Return(mockResult, nil)
			fbAuth.EXPECT().SetCustomClaims(externalUserID, gomock.Any()).Return(nil)
			mock.ExpectCommit()
			Convey("正常に登録できる", func() {
				err := as.Create(mockEntity, password)
				So(err, ShouldBeNil)
			})
		})
		Convey("Update", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}
			mock.ExpectBegin()
			ar.EXPECT().FindByID(gomock.Any(), itemID).Return(mockEntity, nil)
			mock.ExpectCommit()
			Convey("更新対象の管理画面ユーザーが取得できること", func() {
				result, err := as.FindByID(itemID)
				So(result, ShouldResemble, mockEntity)
				So(err, ShouldBeNil)
				mock.ExpectBegin()
				ar.EXPECT().Update(gomock.Any(), itemID, mockEntity).Return(nil)
				mock.ExpectCommit()
				Convey("正常に更新できる", func() {
					err := as.Update(itemID, mockEntity)
					So(err, ShouldBeEmpty)
				})
			})
		})
		Convey("Delete", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}
			Convey("削除対象の管理ユーザーが取得できること", func() {
				mock.ExpectBegin()
				ar.EXPECT().FindByID(gomock.Any(), itemID).Return(mockEntity, nil)
				fbAuth.EXPECT().DeleteUser(externalUserID).Return(nil)
				ar.EXPECT().Delete(gomock.Any(), itemID).Return(nil)
				mock.ExpectCommit()
				Convey("正常に削除できる", func() {
					err := as.Delete(itemID)
					So(err, ShouldBeEmpty)
				})
			})
		})
		Convey("FindAllで正常に取得できなかった場合エラーを返す", func() {
			mock.ExpectBegin()
			ar.EXPECT().FindAll(gomock.Any()).Return(nil, appErr.DBClientError)
			mock.ExpectCommit()

			result, err := as.FindAll()
			So(result, ShouldBeNil)
			So(err, ShouldEqual, appErr.ServiceClientError)
		})
		Convey("Findで正常に取得できなかった場合エラーを返す", func() {
			mockRequest := request.GetItemRequest{
				Name: name, EmailAddress: emailAddress,
			}

			mock.ExpectBegin()
			ar.EXPECT().Find(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, appErr.DBClientError)
			mock.ExpectCommit()

			result, err := as.Find(&mockRequest)
			So(result, ShouldBeNil)
			So(err, ShouldEqual, appErr.ServiceClientError)
		})
		Convey("FindByIDで正常に取得できなかった場合エラーを返す", func() {
			mock.ExpectBegin()
			ar.EXPECT().FindByID(gomock.Any(), gomock.Any()).Return(nil, appErr.DBClientError)
			mock.ExpectCommit()

			result, err := as.FindByID(itemID)
			So(result, ShouldBeNil)
			So(err, ShouldEqual, appErr.ServiceClientError)
		})
		Convey("Createで正常に登録できなかった場合エラーを返す", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}
			mockResult := &auth.UserRecord{
				UserInfo: &auth.UserInfo{
					DisplayName: name, Email: emailAddress, UID: externalUserID,
				},
				CustomClaims:  map[string]interface{}{"role": role},
				Disabled:      false,
				EmailVerified: false,
			}

			Convey("CreateUserでエラーが発生", func() {
				mock.ExpectBegin()
				fbAuth.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil, appErr.FirebaseCreateUserError)
				mock.ExpectCommit()

				err := as.Create(mockEntity, password)
				So(err, ShouldEqual, appErr.ServiceClientError)
			})
			Convey("Createでエラーが発生", func() {
				mock.ExpectBegin()
				fbAuth.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(mockResult, nil)
				ar.EXPECT().Create(gomock.Any(), gomock.Any()).Return(appErr.DBClientError)
				fbAuth.EXPECT().DeleteUser(gomock.Any()).Return(nil)
				mock.ExpectCommit()

				err := as.Create(mockEntity, password)
				So(err, ShouldEqual, appErr.ServiceClientError)
			})
			Convey("SetCustomClaimsでエラーが発生", func() {
				mock.ExpectBegin()
				ar.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
				fbAuth.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(mockResult, nil)
				fbAuth.EXPECT().SetCustomClaims(gomock.Any(), gomock.Any()).Return(appErr.FirebaseSetCustomClaimsError)
				fbAuth.EXPECT().DeleteUser(gomock.Any()).Return(nil)
				mock.ExpectCommit()

				err := as.Create(mockEntity, password)
				So(err, ShouldEqual, appErr.ServiceClientError)
			})
		})
		Convey("Updateで正常に更新できなかった場合エラーを返す", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}

			mock.ExpectBegin()
			ar.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(appErr.DBClientError)
			mock.ExpectCommit()

			err := as.Update(itemID, mockEntity)
			So(err, ShouldEqual, appErr.ServiceClientError)
		})
		Convey("DeleteでFirebaseから削除出来なかった場合にエラーを返す", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}
			mock.ExpectBegin()
			ar.EXPECT().FindByID(gomock.Any(), itemID).Return(mockEntity, nil)
			fbAuth.EXPECT().DeleteUser(externalUserID).Return(appErr.FirebaseDeleteUserError)
			mock.ExpectCommit()
			err := as.Delete(itemID)
			So(err, ShouldEqual, appErr.ServiceStatusBadRequestError)
		})
		Convey("DeleteでDBから削除出来なかった場合にエラーを返す", func() {
			mockEntity := &entities.Item{
				Name: name, ExternalUserID: externalUserID, EmailAddress: emailAddress, Role: role,
			}
			mock.ExpectBegin()
			ar.EXPECT().FindByID(gomock.Any(), itemID).Return(mockEntity, nil)
			fbAuth.EXPECT().DeleteUser(externalUserID).Return(nil)
			ar.EXPECT().Delete(gomock.Any(), itemID).Return(appErr.DBClientError)
			mock.ExpectCommit()
			err := as.Delete(itemID)
			So(err, ShouldEqual, appErr.ServiceClientError)
		})
	})
}

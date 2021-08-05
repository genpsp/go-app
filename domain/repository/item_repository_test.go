package repositories

import (
	entities "github.com/genpsp/go-app/domain/entities"
	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
	"gorm.io/gorm"
	"testing"
)

func TestItemRepositoryImpl_FindAll(t *testing.T) {
	truncateTable("item")
	repository := &ItemRepositoryImpl{}

	id := uint(1)
	name := "name"

	item := entities.Item{
		Model: gorm.Model{
			ID:        id,
			CreatedAt: mock_now,
			UpdatedAt: mock_now,
		},
		Name: name,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("データが存在しない場合空配列を返すこと", t, func() {
		actual, err := repository.FindAll(test_db.Master)
		So(err, ShouldBeNil)
		So(actual, ShouldBeEmpty)
	})

	Convey("データが存在していた場合正しくentityを返すこと", t, func() {
		_ = repository.Create(test_db.Master, &item)

		expect := []entities.Item{{
			Model: gorm.Model{
				ID:        id,
				CreatedAt: mock_now,
				UpdatedAt: mock_now,
			},
			Name: name,
		}}

		actual, err := repository.FindAll(test_db.Master)

		So(err, ShouldBeNil)
		So(actual, ShouldResemble, &expect)
	})
}

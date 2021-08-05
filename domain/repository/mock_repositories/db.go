package mock_repositories

import (
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDBMock() (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gdb, err := gorm.Open(mysql.Dialector{Config: &mysql.Config{DriverName: "mysql",
		Conn: db, SkipInitializeWithVersion: true}},
		&gorm.Config{})

	if err != nil {
		return nil, nil, err
	}
	return gdb, mock, nil
}

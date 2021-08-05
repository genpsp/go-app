package database

import (
	"fmt"
	mysqlcfg "github.com/genpsp/go-app/pkg/configs/mysql"
	"github.com/genpsp/go-app/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

type Database struct {
	Master *gorm.DB
}

func (d *Database) Close() {
	closedDB, _ := d.Master.DB()
	if err := closedDB.Close(); err != nil {
		logger.Logging.Error(fmt.Sprintf("masterDB connection close error. %v", err))
	} else {
		logger.Logging.Info("masterDB connection close success")
	}
}

func dataSource(userName string, password string, host string, dbName string) string {
	return fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local", userName, password, host, dbName)
}

func Open(cfg mysqlcfg.MySql) Database {
	master, err := gorm.Open(mysql.Open(dataSource(cfg.MasterUsername, cfg.MasterPassword, cfg.MasterHost, cfg.DBName)), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{SingularTable: true},
	})
	if err != nil {
		logger.Logging.Fatal(fmt.Sprintf("master connection failed. %v", err))
	}
	if cfg.DebugMode {
		master = master.Debug()
	}

	dbConfig, _ := master.DB()

	dbConfig.SetMaxOpenConns(cfg.MaxOpenConns)
	dbConfig.SetMaxIdleConns(cfg.MaxIdleConns)
	dbConfig.SetConnMaxLifetime(time.Hour)

	logger.Logging.Info(fmt.Sprintf("master connection success. host: %s", cfg.MasterHost))

	return Database{
		Master: master,
	}
}

package mysql

import (
	"fmt"
	"os"

	"github.com/genpsp/go-app/pkg/env"
	"github.com/genpsp/go-app/pkg/utils"
)

type MySql struct {
	MasterUsername   string
	MasterPassword   string
	MasterHost       string
	MasterInstanceID string

	DBName       string
	DebugMode    bool
	MaxOpenConns int
	MaxIdleConns int
}

func NewConfig(env env.Env) MySql {
	var masterHost string
	switch env.ENV {
	case "dev", "stg", "prd":
		masterHost = fmt.Sprintf("unix(/cloudsql/%s)", os.Getenv("CLOUD_SQL_INSTANCE"))
	default:
		masterHost = fmt.Sprintf("tcp(%s)", os.Getenv("MYSQL_MASTER_HOST"))
	}
	return MySql{
		MasterUsername:   env.MasterUsername,
		MasterPassword:   env.MasterPassword,
		MasterHost:       masterHost,
		MasterInstanceID: env.MasterInstanceID,
		DBName:           env.DBName,
		MaxOpenConns:     utils.ConvertInt(env.MaxOpenConns),
		MaxIdleConns:     utils.ConvertInt(env.MaxIdleConns),
		DebugMode:        utils.ConvertBool(env.DebugMode),
	}
}

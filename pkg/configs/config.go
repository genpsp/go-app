package configs

import (
	"github.com/genpsp/go-app/pkg/configs/cloudfunctions"
	"sync"

	"github.com/genpsp/go-app/pkg/configs/firebase"
	"github.com/genpsp/go-app/pkg/configs/gcs"
	"github.com/genpsp/go-app/pkg/configs/logger"
	"github.com/genpsp/go-app/pkg/configs/mysql"
	"github.com/genpsp/go-app/pkg/configs/system"
	env "github.com/genpsp/go-app/pkg/env"
)

var Config *Configuration
var once sync.Once

type Configuration struct {
	MySQL mysql.MySql
}

func LoadConfig() {
	once.Do(func() {
		env := env.NewEnv()

		Config = &Configuration{
			MySQL: mysql.NewConfig(env),
		}
	})
}

func GetConfig() *Configuration {
	return Config
}

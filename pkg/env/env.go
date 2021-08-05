package env

import (
	"os"
)

type Env struct {
	ENV string
}

func NewEnv() Env {
	return Env{
		ENV: os.Getenv("ENV"),
	}
}

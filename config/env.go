package config

import (
	"os"
	"pengoe/internal/logger"

	"github.com/peterszarvas94/envloader"
)

type appConfig struct {
	DB_URL      string
	DB_TOKEN    string
	JWT_SECRET  string
	ENVIRONMENT string
}

func newAppConfig() *appConfig {
	log := logger.Get()

	file, err := os.Open(".env")
	if err == nil {
		envloader.File(file)
	}

	var config appConfig

	err = envloader.Load(&config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &config
}

var Env = newAppConfig()

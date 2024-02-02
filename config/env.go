package config

import (
	"os"
	"path/filepath"
	"pengoe/internal/logger"
	"pengoe/internal/utils"

	"github.com/peterszarvas94/envloader"
)

type appConfig struct {
	DB_URL      string
	DB_TOKEN    string
	JWT_SECRET  string
	ENVIRONMENT string
}

func newAppConfig() *appConfig {
	rootDir, err := utils.GetRootDir()
	if err != nil {
		logger.Fatal(err.Error())
	}

	envFilePath := filepath.Join(rootDir, ".env")

	file, err := os.Open(envFilePath)
	if err == nil {
		envloader.File(file)
	}

	var config appConfig

	err = envloader.Load(&config)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return &config
}

var Env = newAppConfig()

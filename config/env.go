package config

import (
	"fmt"
	"os"
	"path/filepath"
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
	rootDir, err := GetRootDir()
	if err != nil {
		logger.Fatal(err.Error())
	}

	envFilePath := filepath.Join(rootDir, ".env")

	file, err := os.Open(envFilePath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	envloader.File(file)

	var config appConfig

	err = envloader.Load(&config)
	if err != nil {
		logger.Fatal(err.Error())
	}

	return &config
}

var Env = newAppConfig()

// Returns the root directory of the project
func GetRootDir() (string, error) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse upwards until a go.mod file is found
	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		_, err := os.Stat(goModPath)
		if err == nil {
			return currentDir, nil
		}

		// Move one directory up
		parent := filepath.Dir(currentDir)

		// Check if we have reached the root directory
		if parent == currentDir {
			return "", fmt.Errorf("go.mod file not found")
		}

		// Continue the loop with the parent directory
		currentDir = parent
	}
}

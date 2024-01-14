package utils

import (
	"errors"
	"fmt"
	"os"
	"pengoe/internal/logger"
)

func getVariable(key string) (string, error) {
	env, found := os.LookupEnv(key)
	if !found || env == "" {
		return "", errors.New(fmt.Sprintf("Environment variable %s not found", key))
	}

	return env, nil
}

type environmentVariables struct {
	DBUrl       string
	DBToken     string
	JWTSecret   string
	Environment string
}

func initEnvironmentVariables(variables *environmentVariables, keys ...string) {
	for _, key := range keys {
		value, err := getVariable(key)
		if err != nil {
			logger.Log(logger.ERROR, "env/initEnvVars", err.Error())
			os.Exit(1)
		}

		switch key {
		case "DB_URL":
			variables.DBUrl = value
		case "DB_TOKEN":
			variables.DBToken = value
		case "JWT_SECRET":
			variables.JWTSecret = value
		case "ENVIRONMENT":
			variables.Environment = value
		}
	}
}

var Env environmentVariables

func init() {
	initEnvironmentVariables(
		&Env,
		"DB_URL",
		"DB_TOKEN",
		"JWT_SECRET",
		"ENVIRONMENT",
	)
}

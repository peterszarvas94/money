package utils

import (
	"fmt"
	"os"
	"pengoe/internal/logger"
	"reflect"
	"strings"
)

type env struct {
	DB_URL      string
	DB_TOKEN    string
	JWT_SECRET  string
	ENVIRONMENT string
}

func initEnv(variables *env, keys ...string) {
	for _, key := range keys {
		value, found := os.LookupEnv(key)
		if !found {
			logger.Log(
				logger.ERROR,
				"env/initEnvVars",
				fmt.Sprintf("Environment variable %s not found", key),
			)
			os.Exit(1)
		}

		fieldName := strings.ToUpper(key)
		field := reflect.ValueOf(variables).Elem().FieldByName(fieldName)

		if !field.IsValid() {
			logger.Log(logger.ERROR, "env/initEnvVars", fmt.Sprintf("Unknown environment variable: %s", key))
			os.Exit(1)
		}

		if field.Kind() == reflect.String {
			field.SetString(value)
		} else {
			logger.Log(logger.ERROR, "env/initEnvVars", fmt.Sprintf("Unsupported type for field %s", fieldName))
			os.Exit(1)
		}
	}
}

var Env env

func init() {
	initEnv(
		&Env,
		"DB_URL",
		"DB_TOKEN",
		"JWT_SECRET",
		"ENVIRONMENT",
	)
}

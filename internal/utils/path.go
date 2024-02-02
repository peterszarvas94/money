package utils

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func GetPatternFromStr(s string) []string {
	rawPath := strings.Split(s, "/")
	var path []string
	for _, element := range rawPath {
		if element != "" {
			path = append(path, element)
		}
	}
	return path
}

func GetPathVariables(pattern []string, path []string) map[string]string {
	variables := make(map[string]string)

	for i, patternSegment := range pattern {
		if strings.HasPrefix(patternSegment, ":") && i < len(path) {
			variableName := strings.TrimPrefix(patternSegment, ":")
			variables[variableName] = path[i]
		}
	}

	return variables
}

func GetQueryParam(values url.Values, param string) string {
	for key, value := range values {
		if key == param {
			return value[0]
		}
	}
	return ""
}

func IsValidRedirect(redirect string, encoded bool) bool {
	if encoded {
		return strings.HasPrefix(redirect, "%2f") || strings.HasPrefix(redirect, "%2F")
	}

	return strings.HasPrefix(redirect, "/")
}

/*
RemoveTrailingSlash removes trailing slash from path.
*/
func RemoveTrailingSlash(path string) string {
	if path != "/" && strings.HasSuffix(path, "/") {
		return path[:len(path)-1]
	}

	return path
}

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

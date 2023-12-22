package utils

import (
	"net/url"
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

func GetQueryParam(values url.Values, param string) (string) {
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

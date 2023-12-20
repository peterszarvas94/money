package utils

import (
	"net/http"
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

func GetPathVariables(path []string, pattern []string) map[string]string {
	variables := make(map[string]string)

	for i, patternSegment := range pattern {
		if strings.HasPrefix(patternSegment, ":") && i < len(path) {
			variableName := strings.TrimPrefix(patternSegment, ":")
			variables[variableName] = path[i]
		}
	}

	return variables
}

func GetQueryParams(r *http.Request) map[string]string {
	queryParams := make(map[string]string)

	query := r.URL.Query()
	for key, value := range query {
		encoded := url.QueryEscape(value[0])
		queryParams[key] = encoded
	}

	return queryParams
}

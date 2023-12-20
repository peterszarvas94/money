package utils

import (
	"net/http"
	"net/url"
	"strings"
)

/*
GetPath returns a slice of strings from a url path
Example:
urlPath: /users/123
returns: []string{"users", "123"}
*/
// func GetPath(r *http.Request) []string {
// 	rawPath := strings.Split(r.URL.Path, "/")
// 	var path []string
// 	for _, element := range rawPath {
// 		if element != "" {
// 			path = append(path, element)
// 		}
// 	}
// 	return path
// }


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

/*
GetPathVariables returns a map of path variables from a url path and a pattern
Example:
urlPath: /users/123
pattern: /users/:id
returns: map[string]string{"id": "123"}
*/
// func GetPathVariables(path, pattern string) map[string]string {
// 	variables := make(map[string]string)
//
// 	urlParts := strings.Split(path, "/")
// 	patternParts := strings.Split(pattern, "/")
//
// 	for i, part := range patternParts {
// 		if strings.HasPrefix(part, ":") && i < len(urlParts) {
// 			variableName := strings.TrimPrefix(part, ":")
// 			variables[variableName] = urlParts[i]
// 		}
// 	}
//
// 	return variables
// }

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

/*
GetQueryParams returns a map of query parameters from a url query
Example:
urlQuery: ?id=123&name=John
returns: map[string]string{"id": "123", "name": "John"}
*/
func GetQueryParams(r *http.Request) map[string]string {
	queryParams := make(map[string]string)

	query := r.URL.Query()
	for key, value := range query {
		encoded := url.QueryEscape(value[0])
		queryParams[key] = encoded
	}

	return queryParams
}

package utils

import (
	"errors"
	"fmt"
	"html"
	"net/url"
)

func GetFormValues(values url.Values, fields ...string) (map[string]string, error) {
	result := make(map[string]string)

	for _, field := range fields {
		value := html.EscapeString(values.Get(field))
		if value == "" {
			return nil, errors.New(fmt.Sprintf("Form value \"%s\" is empty", field))
		}
		result[field] = value
	}

	return result, nil
}

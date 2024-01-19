package utils_test

import (
	"net/url"
	"pengoe/internal/utils"
	"testing"
)

func TestGetFormValues(t *testing.T) {
	var values url.Values = make(map[string][]string)
	values.Set("username", "test")
	values.Set("email", "")

	_, err := utils.GetFormValues(values, "username", "email")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	result2, err2 := utils.GetFormValues(values, "username")
	if err2 != nil {
		t.Errorf("Expected nil, got %s", err2.Error())
	}
	if result2["username"] != "test" {
		t.Errorf("Expected %s, got %s", "test", result2["username"])
	}

	result3, err3 := utils.GetFormValues(values, "username", "idk")
	if err3 == nil {
		t.Errorf("Expected error, got nil")
	}
	if result3 != nil {
		t.Errorf("Expected nil, got %s", result3)
	}
}

package utils_test

import (
	"net/url"
	"pengoe/internal/utils"
	"testing"
)

func TestGetPatternFromString(t *testing.T) {
	s := "/test/path"
	expected := []string{"test", "path"}

	result := utils.GetPatternFromStr(s)

	if !utils.SliceEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetPathVariables(t *testing.T) {
	pattern := []string{"test", ":id", "path"}
	path := []string{"test", "123", "path"}
	expected := map[string]string{"id": "123"}

	result := utils.GetPathVariables(pattern, path)

	if !utils.MapEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetQueryParams(t *testing.T) {
	var values url.Values
	values = make(map[string][]string)
	values["test"] = []string{"123", "456"}
	expected := "123"

	result:= utils.GetQueryParam(values, "test")
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestIsValidEncodedRedirect(t *testing.T) {
	redirect1 := "%2ftest"
	expected1 := true
	result1 := utils.IsValidRedirect(redirect1, true)
	if result1 != expected1 {
		t.Errorf("Expected1 %v, got %v", expected1, result1)
	}

	redirect2 := "https%3A%2F%2Fexample.com"
	expected2 := false
	result2 := utils.IsValidRedirect(redirect2, true)
	if result2 != expected2 {
		t.Errorf("Expected2 %v, got %v", expected2, result2)
	}

	redirect3 := "/test"
	expected3 := false
	result3 := utils.IsValidRedirect(redirect3, true)
	if result3 != expected3 {
		t.Errorf("Expected3 %v, got %v", expected3, result3)
	}

	redirect4 := "http://example.com/test"
	expected4 := false
	result4 := utils.IsValidRedirect(redirect4, true)
	if result4 != expected4 {
		t.Errorf("Expected4 %v, got %v", expected4, result4)
	}
}

func TestIsValidDecodedRedirect(t *testing.T) {
	redirect1 := "%2ftest"
	expected1 := false
	result1 := utils.IsValidRedirect(redirect1, false)
	if result1 != expected1 {
		t.Errorf("Expected1 %v, got %v", expected1, result1)
	}

	redirect2 := "https%3A%2F%2Fexample.com"
	expected2 := false
	result2 := utils.IsValidRedirect(redirect2, false)
	if result2 != expected2 {
		t.Errorf("Expected2 %v, got %v", expected2, result2)
	}

	redirect3 := "/test"
	expected3 := true
	result3 := utils.IsValidRedirect(redirect3, false)
	if result3 != expected3 {
		t.Errorf("Expected3 %v, got %v", expected3, result3)
	}

	redirect4 := "http://example.com/test"
	expected4 := false
	result4 := utils.IsValidRedirect(redirect4, false)
	if result4 != expected4 {
		t.Errorf("Expected4 %v, got %v", expected4, result4)
	}
}

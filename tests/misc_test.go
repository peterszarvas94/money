package tests

import (
	"pengoe/internal/utils"
	"testing"
)

func TestSliceEqual(t *testing.T) {
	a1 := []string{"test", "helo"}
	b1 := []string{"test", "helo"}
	expected1 := true

	result1 := utils.SliceEqual(a1, b1)
	if result1 != expected1 {
		t.Errorf("Expected1 %v, got %v", expected1, result1)
	}

	a2 := []string{"test", "helo"}
	b2 := []string{"test", "helo", "other"}
	expected2 := false
	if utils.SliceEqual(a2, b2) != expected2 {
		t.Errorf("Expected2 %v, got %v", expected2, utils.SliceEqual(a2, b2))
	}
}

func TestGetFromSlice(t *testing.T) {
	slice := []string{"test", "helo"}
	expected := "test"

	result, err := utils.GetFromSlice(0, slice)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}

	_, err = utils.GetFromSlice(2, slice)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestMapEqual(t *testing.T) {
	a1 := map[string]string{"test": "helo"}
	b1 := map[string]string{"test": "helo"}
	expected1 := true

	result1 := utils.MapEqual(a1, b1)
	if result1 != expected1 {
		t.Errorf("Expected1 %v, got %v", expected1, result1)
	}

	a2 := map[string]string{"test": "helo"}
	b2 := map[string]string{"test": "helo", "other": "test"}
	expected2 := false
	if utils.MapEqual(a2, b2) != expected2 {
		t.Errorf("Expected2 %v, got %v", expected2, utils.MapEqual(a2, b2))
	}
}

package utils

import (
	"testing"
	"time"
)

func TestConvertToTime(t *testing.T) {
	expected := time.Date(2026, 1, 2, 15, 4, 5, 0, time.UTC)
	expectedStr := expected.Format(time.RFC3339Nano)

	converted, err := ConvertToTime(expectedStr)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if converted != expected {
		t.Errorf("Expected %v, got %v", expected, converted)
	}
}

package utils

import (
	"testing"
)

func TestHash(t *testing.T) {
	password := "test"
	hashedPassword, hashErr := HashPassword(password)
	if hashErr != nil {
		t.Errorf("Expected no error, got %v", hashErr)
	}

	matchErr := CheckPasswordHash(hashedPassword, password)
	if matchErr != nil {
		t.Errorf("Expected no error, got %v", matchErr)
	}
}

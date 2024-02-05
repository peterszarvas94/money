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

	match := CheckPasswordHash(hashedPassword, password)
	if !match {
		t.Errorf("Expected match to be true, got false")
	}
}

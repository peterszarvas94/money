package utils_test

import (
	"pengoe/internal/utils"
	"testing"
)

func TestHash(t *testing.T) {
	password := "test"
	hashedPassword, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		t.Errorf("Expected no error, got %v", hashErr)
	}

	matchErr := utils.CheckPasswordHash(hashedPassword, password)
	if matchErr != nil {
		t.Errorf("Expected no error, got %v", matchErr)
	}
}

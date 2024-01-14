package utils_test

import (
	"pengoe/internal/utils"
	"testing"
)

func TestGenerateCSRFToken(t *testing.T) {
	token, err := utils.GenerateCSRFToken()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(token) != 44 {
		t.Errorf("Expected token to be 44 characters long, got %d", len(token))
	}
}

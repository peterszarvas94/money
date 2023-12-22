package utils

import (
	"strconv"
	"testing"
	"time"
)

func TestNewToken(t *testing.T) {
	id := 1
	variant := AccessToken

	token, err := NewToken(id, variant)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if token.Expires <= time.Now().Unix() {
		t.Errorf("Expected token to expire in the future, got %d", token.Expires)
	}

	validatedToken, validationErr := ValidateToken(token.Token)
	if validationErr != nil {
		t.Errorf("Expected no error, got %v", validationErr)
	}

	subjec, subjectErr := validatedToken.GetSubject()
	if subjectErr != nil {
		t.Errorf("Expected no error, got %v", subjectErr)
	}

	subjectInt, converErr := strconv.Atoi(subjec)
	if converErr != nil {
		t.Errorf("Expected no error, got %v", converErr)
	}

	if subjectInt != id {
		t.Errorf("Expected subject to be %d, got %d", id, subjectInt)
	}
}

package utils

import (
	"strings"
	"testing"
)

func TestGenerateUUID(t *testing.T) {
	uuid := NewUUID("test")
	if len(uuid) == 0 {
		t.Errorf("UUID is empty")
	}

	withoutPrefix := strings.Split(uuid, "_")[1]
	if len(withoutPrefix) != 32 {
		t.Errorf("UUID is not 32 characters long")
	}
}

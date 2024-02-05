package utils

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

/*
NewUUID is a function that returns a UUID with a prefix.
Eg. test_1234567890abcdef1234567890abcdef
*/
func NewUUID(prefix string) string {
	str := uuid.New().String()
	clean := strings.ReplaceAll(str, "-", "")
	return fmt.Sprintf("%s_%s", prefix, clean)
}

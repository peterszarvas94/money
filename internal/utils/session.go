package utils

import (
	"crypto/rand"
	"encoding/base64"
)

/*
generateCSRFToken generates a random 32 byte string and encodes it to base64.
*/
func GenerateCSRFToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

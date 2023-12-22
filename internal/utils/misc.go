package utils

import "errors"

func SliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, elem := range a {
		if elem != b[i] {
			return false
		}
	}
	return true
}

func GetFromSlice(i int, s []string) (string, error) {
	if i < len(s) {
		return s[i], nil
	}
	return "", errors.New("Index out of range")
}

func MapEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for key, value := range a {
		if b[key] != value {
			return false
		}
	}
	return true
}

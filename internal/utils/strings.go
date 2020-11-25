package utils

import (
	"crypto/rand"
	"strings"
)

// GenerateRandomString returns a crypto-strong random string of specified length.
func GenerateRandomString(n int) (string, error) {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

	// generate a random slice of bytes
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = alphabet[b%byte(len(alphabet))]
	}
	return string(bytes), nil
}

func ContainsString(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func RemoveString(arr []string, str string) []string {
	for i, v := range arr {
		if v == str {
			return append(arr[:i], arr[i+1:]...)
		}
	}
	return arr
}

func TrimQuotes(s string) string {
	if strings.HasPrefix(s, "\"") {
		return strings.Trim(s, "\"")
	} else if strings.HasPrefix(s, "'") {
		return strings.Trim(s, "'")
	} else {
		return s
	}
}

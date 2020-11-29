package utils

import (
	"crypto/rand"
	"encoding/json"
	"regexp"
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

// https://stackoverflow.com/questions/46128016/insert-a-value-in-a-slice-at-a-given-index
func InsertString(arr []string, value string, index int) []string {
	if index >= len(arr) { // nil or empty slice or after last element
		return append(arr, value)
	}
	arr = append(arr[:index+1], arr[index:]...) // index < len(a)
	arr[index] = value
	return arr
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

func ReplaceEscapeSymbols(s string) string {
	r := strings.ReplaceAll(s, "\\n", "\n")
	r = strings.ReplaceAll(r, "\\t", "\t")
	r = strings.ReplaceAll(r, "\\'", "'")
	r = strings.ReplaceAll(r, "\\\"", "\"")
	return r
}

func WhiteSpaceOnly(s string) bool {
	re := regexp.MustCompile(`[^\s]+`)
	return !re.Match([]byte(s))
}

func CountIndent(s string) int {
	re := regexp.MustCompile(`\s*`)
	r := string(re.Find([]byte(s)))
	r = strings.ReplaceAll(r, "\t", "    ")
	return len(r)
}

func CountRune(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}

func ToJSON(v interface{}) string {
	json, _ := json.MarshalIndent(v, "", "  ")
	return string(json)
}

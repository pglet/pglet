package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
)

// SHA1 returns SHA1 hash of the input string.
func SHA1(value string) string {
	h := sha1.New()
	io.WriteString(h, value)
	return fmt.Sprintf("%x", h.Sum(nil))
}

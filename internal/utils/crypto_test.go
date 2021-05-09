package utils

import (
	"testing"
)

func TestShortGetCipherSecretKey(t *testing.T) {
	shortKey := "123"
	chipherKey := getCipherKey(shortKey)

	l := len(chipherKey)
	if l != 32 {
		t.Errorf("getCipherSecretKey returned array of %d bytes length, want %d", l, 32)
	}

	if chipherKey[3] != 0 && chipherKey[l-1] != 0 {
		t.Error("For short key 4th and last bytes must be zeroes")
	}
}

func TestLongGetCipherSecretKey(t *testing.T) {
	longKey := "0123456789abcdefghijkl0123456789-somemoresymbols"
	chipherKey := getCipherKey(longKey)

	l := len(chipherKey)
	if l != 32 {
		t.Errorf("getCipherSecretKey returned array of %d bytes length, want %d", l, 32)
	}

	if chipherKey[3] != longKey[3] && chipherKey[l-1] != longKey[l-1] {
		t.Error("Generated chipher key is not valid")
	}
}

func TestEncryptWithKey(t *testing.T) {

	// test data
	secretKey := "secret!"
	plaintext := "Hello, world!"

	// encrypt
	encData, err := EncryptWithKey([]byte(plaintext), secretKey)
	if err != nil {
		//t.Errorf("TestUnquote3 returned %s, want %s", r, expR)
		t.Error("Encryption should not fail")
	}

	// decrypt
	decData, err := DecryptWithKey(encData, secretKey)
	if err != nil {
		t.Error("Decryption should not fail")
	}

	if plaintext != string(decData) {
		t.Errorf("DecryptWithKey returned %s, want %s", decData, plaintext)
	}
}

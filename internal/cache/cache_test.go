package cache

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// test Redis
	os.Setenv("REDIS_ADDR", "localhost:6379")
	Init()
	retCode := m.Run()
	if retCode != 0 {
		os.Exit(retCode)
	}

	// test in-memory
	os.Setenv("REDIS_ADDR", "")
	Init()
	retCode = m.Run()
	if retCode != 0 {
		os.Exit(retCode)
	}
}

func TestGetString(t *testing.T) {
	v := "111"
	SetString("aaa", v, 0)
	r := GetString("aaa")
	if r != v {
		t.Errorf("getString returned %s, want %s", r, v)
	}
}

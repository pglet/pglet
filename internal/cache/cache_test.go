package cache

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// test Redis
	// os.Setenv("REDIS_ADDR", "localhost:6379")
	// Init()
	// retCode := m.Run()
	// if retCode != 0 {
	// 	os.Exit(retCode)
	// }

	// test in-memory
	os.Setenv("REDIS_ADDR", "")
	Init()
	retCode := m.Run()
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

func TestInc(t *testing.T) {

	c := Inc("inc1", 1)
	if c != 1 {
		t.Errorf("inc returned %d, want %d", c, 1)
	}
	c = Inc("inc1", 2)
	if c != 3 {
		t.Errorf("inc returned %d, want %d", c, 3)
	}
}

func TestHash(t *testing.T) {

	n0 := HashGet("non-existent hash", "something")
	if n0 != "" {
		t.Errorf("HashGet of non-existent key returned %s, want %s", n0, "")
	}

	HashSet("hash1", "aaa", "1", "bbb", "Test")
	aaa := HashGet("hash1", "aaa")
	if aaa != "1" {
		t.Errorf("HashGet returned %s, want %s", aaa, "1")
	}
	bbb := HashGet("hash1", "bbb")
	if bbb != "Test" {
		t.Errorf("HashGet returned %s, want %s", bbb, "Test")
	}
	n1 := HashGet("hash1", "something")
	if n1 != "" {
		t.Errorf("HashGet non-existent field returned %s, want %s", n1, "")
	}

	HashSet("hash1", "ccc", "Another test")

	entries := HashGetAll("hash1")
	count := len(entries)
	if count != 3 {
		t.Errorf("HashGetAll returned %d entries, want %d", count, 3)
	}

	e1 := entries["aaa"]
	if e1 != "1" {
		t.Errorf("Checking all entries field 'aaa' returned %s, want %s", e1, "1")
	}
	e2 := entries["ccc"]
	if e2 != "Another test" {
		t.Errorf("Checking all entries field 'ccc' returned %s, want %s", e2, "Another test")
	}

	HashRemove("hash1", "aaa")
	HashRemove("hash1", "bbb")
	entries = HashGetAll("hash1")
	count = len(entries)
	if count != 1 {
		t.Errorf("HashGetAll after removing 2 elements returned %d entries, want %d", count, 1)
	}
	HashRemove("hash1", "ccc")
	if Exists("hash1") {
		t.Errorf("Hash should not exist after all its elements deleted")
	}
}

func TestSet(t *testing.T) {
	if Exists("set1") {
		t.Errorf("Set should not exist in the first place")
	}

	SetAdd("set1", "v1")
	SetAdd("set1", "v1")
	items := SetGet("set1")
	count := len(items)
	if count != 1 {
		t.Errorf("SetGet returned %d entries, want %d", count, 1)
	}

	SetAdd("set1", "v2")
	items = SetGet("set1")
	count = len(items)
	if count != 2 {
		t.Errorf("SetGet returned %d entries, want %d", count, 2)
	}

	SetRemove("set1", "v2")
	items = SetGet("set1")
	count = len(items)
	if count != 1 {
		t.Errorf("SetGet after removing v2 returned %d entries, want %d", count, 1)
	}

	SetRemove("set1", "v1")
	if Exists("set1") {
		t.Errorf("Set should not exist after removing all elements")
	}
}

// func TestChannels(t *testing.T) {
// 	ch1 := make(chan int)
// 	log.Println(ch1)

// 	go func() {
// 		for {
// 			select {
// 			case msg, more := <-ch1:
// 				if !more {
// 					fmt.Println("Channel closed")
// 					return
// 				}
// 				log.Println(msg)
// 			}
// 		}
// 	}()

// 	ch1 <- 1
// 	ch1 <- 2
// 	ch1 <- 3
// 	close(ch1)
// 	time.Sleep(3 * time.Second)
// }

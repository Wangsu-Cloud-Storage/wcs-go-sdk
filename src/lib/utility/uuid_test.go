package utility

import (
	"testing"
)

func TestGetUuid(t *testing.T) {
	uuid1 := GetUuid()
	uuid2 := GetUuid()
	if uuid1 == uuid2 {
		t.Fatal("Generate 2 of the same UUIDs!")
	}
	t.Log(uuid1)
	t.Log(uuid2)
}

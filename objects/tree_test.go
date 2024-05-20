package objects

import (
	"testing"
)

func TestTreeNames(t *testing.T) {
	objID := "15c570b70d51c18c85e580a4d00c8ea02f1ea2ef"
	tree, err := ReadObjectFromFile(objID)
	if err != nil {
		t.Error(err)
	}

	if tree.HashString() != objID {
		t.Errorf("invalid hash; have %s, want %s", tree.HashString(), objID)
	}
}

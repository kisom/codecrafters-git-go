package objects

import (
	"testing"
)

func TestTreeNames(t *testing.T) {
	objID := "15c570b70d51c18c85e580a4d00c8ea02f1ea2ef"
	blob, err := ReadBlobWithID(objID)
	if err != nil {
		t.Error(err)
	}

	if blob.HashString() != objID {
		t.Errorf("invalid blob hash; have %s, want %s", blob.HashString(), objID)
	}

	tree, err := ReadObjectFromFile(objID)
	if err != nil {
		t.Error(err)
	}

	if tree.HashString() != objID {
		t.Errorf("invalid tree hash; have %s, want %s", tree.HashString(), objID)
	}
}

func TestListTree(t *testing.T) {
	objID := "74633606588dc5375716604aed882c6b106e4530"
	ListTree([]string{objID})
}

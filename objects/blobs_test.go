package objects

import (
	"bytes"
	"testing"
)

func TestReadBlob(t *testing.T) {
	objectID := "176a458f94e0ea5272ce67c36bf30b6be9caf623"
	objectContents := "blob 12\x00* text=auto\n"

	blob, err := readBlob(objectID)
	if err != nil {
		t.Fatalf("failed to read blob: %v", err)
	}

	raw := blob.Raw()
	if !bytes.Equal(raw, []byte(objectContents)) {
		t.Fatalf("blob contents don't match\nexpected: %v\nactual: %v", objectContents, raw)
	}
}

func TestHashBlob(t *testing.T) {
	objectID := "3b18e512dba79e4c8300dd08aeb37f8e728b8dad"
	contents := []byte("hello world\n")

	blob := BlobFromBytes(contents)
	if blob.HashString() != objectID {
		t.Fatalf("blob hash doesn't match\nexpected: %v\nactual: %v", objectID, blob.HashString())
	}
}

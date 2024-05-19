package catfile

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

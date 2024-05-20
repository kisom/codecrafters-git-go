package objects

import (
	"fmt"
	"testing"
)

func TestReadObjectFromFile(t *testing.T) {
	var testCases = []struct {
		ID   string
		Type string
	}{
		{"176a458f94e0ea5272ce67c36bf30b6be9caf623", "blob"},
		{"f7ba871793d3b38d81e217caec01558426c730c5", "tree"},
	}

	for _, tc := range testCases {
		object, err := ReadObjectFromFile(tc.ID)
		if err != nil {
			t.Error(err)
		}

		switch object.(type) {
		case *Blob:
			if tc.Type != "blob" {
				t.Errorf("expected %s, got %T", tc.Type, object)
			}
		case *Tree:
			if tc.Type != "tree" {
				t.Errorf("expected %s, got %T", tc.Type, object)
			}
		default:
			panic(fmt.Sprintf("unknown object type %T", object))
		}
	}
}

func TestReadBasicObjectFromFile(t *testing.T) {
	var testCases = []struct {
		ID   string
		Type string
	}{
		{"176a458f94e0ea5272ce67c36bf30b6be9caf623", "blob"},
		{"f7ba871793d3b38d81e217caec01558426c730c5", "tree"},
	}

	for _, tc := range testCases {
		object, err := ReadBlobWithID(tc.ID)
		if err != nil {
			t.Error(err)
		}

		if object.Type != tc.Type {
			t.Errorf("expected %s, got %s", tc.Type, object.Type)
		}
	}
}

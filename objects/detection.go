package objects

import "github.com/pkg/errors"

type Object interface {
	Raw() []byte
	Hash() []byte
	HashString() string
	String() string
	ObjectType() string
}

const (
	TypeBlob = "blob"
	TypeTree = "tree"
)

func ReadObjectFromFile(id string) (Object, error) {
	obj, err := ReadBlobWithID(id)
	if err != nil {
		return nil, errors.Wrap(err, "reading blob with id "+id+" from file")
	}

	switch obj.Type {
	case "blob":
		return obj, nil
	case "tree":
		return TreeFromBlob(obj)
	}

	panic("unknown object type " + obj.Type)
}

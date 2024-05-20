package objects

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"flag"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/die"
	"github.com/kisom/codecrafters/git-go/kgit"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

var _ Object = &Blob{}

type Blob struct {
	Type     string
	Contents []byte
}

func (blob *Blob) Header() []byte {
	return []byte(fmt.Sprintf("%s %d", blob.Type, blob.Size()))
}

func (blob *Blob) Raw() []byte {
	raw := append(blob.Header(), 0)
	raw = append(raw, blob.Contents...)
	return raw
}

func (blob *Blob) Hash() []byte {
	hash := sha1.Sum(blob.Raw())
	return hash[:]
}

func (blob *Blob) HashString() string {
	return fmt.Sprintf("%x", blob.Hash())
}

func (blob *Blob) Size() int {
	return len(blob.Contents)
}

func (blob *Blob) String() string {
	return string(blob.Contents)
}

func (blob *Blob) ObjectType() string {
	return blob.Type
}

func (blob *Blob) Write() error {
	id := blob.HashString()
	path, err := kgit.PathFromID(id)
	if err != nil {
		return errors.Wrap(err, "couldn't find git path")
	}

	parent := filepath.Dir(path)
	err = os.MkdirAll(parent, 0755)
	if err != nil {
		return errors.Wrap(err, "couldn't create parent directory "+parent)
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "couldn't create file "+path)
	}
	defer file.Close()

	encoder := zlib.NewWriter(file)
	defer encoder.Close()

	_, err = io.Copy(encoder, bytes.NewReader(blob.Raw()))
	if err != nil {
		return errors.Wrap(err, "couldn't write to file "+path)
	}

	return nil
}

func ReadBlobWithID(id string) (*Blob, error) {
	path, err := kgit.PathFromID(id)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get path from id "+id)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open file with id "+id)
	}

	defer file.Close()

	decoder, err := zlib.NewReader(file)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create zlib reader")
	}
	defer decoder.Close()

	contents, err := io.ReadAll(decoder)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read object content for id "+id)
	}

	obj := &Blob{}
	var header []byte

	for i := range contents {
		if contents[i] == 0x0 {
			header = contents[:i]
			obj.Contents = contents[i+1:]
			break
		}
	}

	headerParts := bytes.SplitN(header, []byte(" "), 2)
	if len(headerParts) != 2 {
		fmt.Fprintf(os.Stderr, "header has %d parts", len(headerParts))
		return nil, errors.New("invalid object header for id " + id + ", header: " + string(header))
	}

	obj.Type = string(headerParts[0])
	size, err := strconv.Atoi(string(headerParts[1]))
	if err != nil {
		return nil, errors.Wrap(err, "invalid object header size for id "+id+", header: "+string(header))
	}

	if size != len(obj.Contents) {
		return nil, fmt.Errorf("objects: header size mismatch: %d != %d", obj.Size(), len(obj.Contents))
	}

	return obj, nil
}

func NewBlobFromFile(path string) (*Blob, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't open file at  "+path)
	}
	defer file.Close()

	contents, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't read file at "+path)
	}

	obj := &Blob{
		Type:     "blob",
		Contents: contents,
	}

	return obj, nil
}

func BlobFromBytes(content []byte) *Blob {
	return &Blob{
		Type:     "blob",
		Contents: content[:],
	}
}

func catBlob(id string) {
	object, err := ReadBlobWithID(id)
	die.If(err)

	fmt.Printf("%s", string(object.Contents))
}

func CatFile(args []string) {
	var objectPath string

	flagset := flag.NewFlagSet("cat-file", flag.ExitOnError)
	flagset.StringVar(&objectPath, "p", "", "object `ID` as a hash")
	err := flagset.Parse(args)
	if err != nil {
		die.If(err)
	}

	if objectPath != "" {
		catBlob(objectPath)
	}
}

func Hash(args []string) {
	var writeObjects bool

	flagset := flag.NewFlagSet("hash-object", flag.ExitOnError)
	flagset.BoolVar(&writeObjects, "w", false, "write objects to git repository")
	err := flagset.Parse(args)
	die.If(err)

	succeeded := true
	for _, arg := range flagset.Args() {
		blob, err := NewBlobFromFile(arg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "objects: couldn't read source %s: %v\n", arg, err)
			succeeded = false
			continue
		}

		fmt.Printf("%s\n", blob.HashString())
		if writeObjects {
			err = blob.Write()
			if err != nil {
				fmt.Fprintf(os.Stderr, "objects: couldn't write blob %s with id %s: %v\n",
					arg, blob.HashString(), err)
				succeeded = false
				continue
			}
		}
	}

	if !succeeded {
		os.Exit(1)
	}
}

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

type Blob struct {
	Size     int
	Contents []byte
	hash     []byte
}

func BlobFromFile(path string) (*Blob, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("objects: opening file %s", path))
	}
	defer file.Close()

	blob := &Blob{}
	blob.Contents, err = io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("objects: reading file %s", path))
	}

	blob.Size = len(blob.Contents)
	return blob, nil
}

func BlobFromBytes(bytes []byte) *Blob {
	blob := &Blob{}
	blob.Contents = bytes
	blob.Size = len(bytes)
	return blob
}

func (b *Blob) String() string {
	return string(b.Contents)
}

func (b *Blob) Header() []byte {
	return []byte(fmt.Sprintf("blob %d", b.Size))
}

func (b *Blob) Raw() []byte {
	raw := b.Header()
	raw = append(raw, 0)
	raw = append(raw, b.Contents...)
	return raw
}

func (b *Blob) Hash() []byte {
	if b.hash == nil {
		h := sha1.Sum(b.Raw())
		b.hash = h[:]
	}

	return b.hash
}

func (b *Blob) Decode(r io.Reader) error {
	decoder, err := zlib.NewReader(r)
	if err != nil {
		return errors.Wrap(err, "objects: couldn't create zlib decoder")
	}
	defer decoder.Close()

	var header []byte
	for {
		b := [1]byte{0}
		_, err = decoder.Read(b[:])
		if err == io.EOF {
			break
		}

		if b[0] == 0 {
			break
		}

		header = append(header, b[0])
	}

	if !bytes.HasPrefix(header, []byte("blob ")) {
		return errors.New("objects: blob has an invalid header")
	}

	size, err := strconv.Atoi(string(header[5:]))
	if err != nil {
		return fmt.Errorf("objects: blob has an unreadable size: %#v", header[5:])
	}

	data := make([]byte, size)
	n, err := decoder.Read(data)
	if err != nil && err != io.EOF {
		return errors.Wrap(err, "objects: failed to decompress blob")
	}

	if n != size {
		return fmt.Errorf("objects: blob has an invalid size: %d, but want %d", n, size)
	}

	b.Size = size
	b.Contents = data

	return nil
}

func (b *Blob) Encode(w io.Writer) error {
	encoder := zlib.NewWriter(w)
	defer encoder.Close()

	_, err := encoder.Write(b.Raw())
	if err != nil {
		return errors.Wrap(err, "objects: failed to write blob")
	}

	return nil
}

func (b *Blob) Write() error {
	path, err := kgit.PathFromID(b.HashString())
	if err != nil {
		return errors.Wrap(err, "objects: failed to build path")
	}

	parent := filepath.Dir(path)
	err = os.MkdirAll(parent, 0644)
	if err != nil {
		return errors.Wrap(err, "objects: failed to create directory "+parent)
	}

	file, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "objects: failed to open file "+path)
	}
	defer file.Close()

	return b.Encode(file)
}

func (b *Blob) HashString() string {
	return fmt.Sprintf("%x", b.Hash())
}

func readBlob(id string) (*Blob, error) {
	path, err := kgit.PathFromID(id)
	if err != nil {
		return nil, errors.Wrap(err, "objects: while trying to read blob "+id)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	blob := &Blob{}
	err = blob.Decode(file)
	if err != nil {
		return nil, errors.Wrap(err, "objects: while trying to read blob "+id)
	}

	return blob, nil
}

func catBlob(id string) {
	blob, err := readBlob(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cat-file: couldn't read blob %s: %v", id, err)
		os.Exit(1)
	}

	fmt.Printf("%s", blob)
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
		blob, err := BlobFromFile(arg)
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

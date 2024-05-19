package catfile

import (
	"bytes"
	"compress/zlib"
	"flag"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/die"
	"github.com/kisom/codecrafters/git-go/kgit"
	"github.com/pkg/errors"
	"io"
	"os"
	"strconv"
)

type Blob struct {
	Size     int
	Header   []byte
	Contents []byte
}

func (b *Blob) String() string {
	return string(b.Contents)
}

func (b *Blob) Raw() []byte {
	raw := b.Header
	raw = append(raw, 0)
	raw = append(raw, b.Contents...)
	return raw
}

func readBlob(id string) (*Blob, error) {
	path, err := kgit.PathFromID(id)
	if err != nil {
		return nil, errors.Wrap(err, "catfile: while trying to read blob "+id)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder, err := zlib.NewReader(file)
	if err != nil {
		return nil, errors.Wrap(err, "cat-file: couldn't create zlib decoder")
	}
	defer decoder.Close()

	header := []byte{}
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
		return nil, fmt.Errorf("cat-file: blob %s has an invalid header", id)
	}

	size, err := strconv.Atoi(string(header[5:]))
	if err != nil {
		return nil, fmt.Errorf("cat-file: blob %s has an unreadable size: %#v", id, header[5:])
	}

	data := make([]byte, size)
	n, err := decoder.Read(data)
	if err != nil && err != io.EOF {
		return nil, errors.Wrap(err, "failed to decompress blob")
	}

	if n != size {
		return nil, fmt.Errorf("cat-file: blob %s has an invalid size: %d, but want %d", id, n, size)
	}

	return &Blob{size, header, data}, nil
}

func catBlob(id string) {
	blob, err := readBlob(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cat-file: couldn't read blob %s: %v", id, err)
		os.Exit(1)
	}

	fmt.Printf("%s", blob)
}

func Run(args []string) {
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

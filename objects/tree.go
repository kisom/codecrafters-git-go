package objects

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/die"
	"github.com/kisom/codecrafters/git-go/kgit"
	"os"
	"strings"
)

type entryMode int

const (
	modeRegular entryMode = (iota + 1) << 1
	modeExecutable
	modeSymbolic
	modeDirectory
)

var modeToString = map[entryMode]string{
	modeRegular:    "100644",
	modeExecutable: "100755",
	modeSymbolic:   "120000",
	modeDirectory:  "040000",
}

var modeFromString = map[string]entryMode{
	"100644": modeRegular,
	"100755": modeExecutable,
	"120000": modeSymbolic,
	"040000": modeDirectory,
}

type TreeEntry struct {
	Hash []byte
	Name string
	Mode entryMode
}

func (e *TreeEntry) Raw() []byte {
	// <mode> <name>\0<20_byte_sha>
	raw := append([]byte(fmt.Sprintf("%s %s\x00", modeToString[e.Mode], e.Name)), e.Hash...)
	return raw
}

func (e *TreeEntry) String() string {
	id := fmt.Sprintf("%0x", e.Hash)
	obj, err := ReadBlobWithID(id)
	if err != nil {
		panic(fmt.Sprintf("object %s can't be read: %v", id, err))
	}

	return fmt.Sprintf("%s %s %s\t%s", modeToString[e.Mode], obj.Type, id, e.Name)
}

func entryFromBytes(a, b []byte) *TreeEntry {
	entry := &TreeEntry{}
	s := strings.Fields(string(a))
	entry.Mode = modeFromString[s[0]]
	entry.Name = s[1]
	entry.Hash = b
	return entry
}

type Tree struct {
	entries []*TreeEntry
}

func (tree *Tree) Hash() []byte {
	hash := sha1.Sum(tree.Raw())
	return hash[:]
}

func (tree *Tree) HashString() string {
	return fmt.Sprintf("%0x", tree.Hash())
}

func (tree *Tree) String() string {
	s := ""
	for _, entry := range tree.entries {
		s += entry.String() + "\n"
	}
	return s
}

func (tree *Tree) rawEntries() []byte {
	var raw []byte
	for _, entry := range tree.entries {
		raw = append(raw, entry.Raw()...)
	}
	return raw
}

func (tree *Tree) Size() int {
	return len(tree.rawEntries())
}

func (tree *Tree) Raw() []byte {
	raw := []byte(fmt.Sprintf("tree %d\x00", tree.Size()))
	raw = append(raw, tree.rawEntries()...)
	return raw
}

func (tree *Tree) Names() string {
	s := ""
	for _, entry := range tree.entries {
		s += entry.Name + "\n"
	}
	return s
}

func (tree *Tree) ObjectType() string {
	return TypeTree
}

func scanEntries(contents []byte) [][]byte {
	lines := [][]byte{}
	line := []byte{}

	for i := 0; i < len(contents); i++ {
		if contents[i] != '\x00' {
			line = append(line, contents[i])
			continue
		}

		lines = append(lines, line)
		i++
		line = []byte{}
		if len(contents[i:]) < 20 {
			panic("invalid tree has incomplete hash")
		}

		line = append(line, contents[i:i+20]...)
		i += 20
	}

	if len(line) > 0 {
		lines = append(lines, line)
	}
	return lines
}

func TreeFromBlob(blob *Blob) (*Tree, error) {
	tree := &Tree{}

	entryData := kgit.Scanner(blob.Contents)
	for i := 0; i < len(entryData); i += 2 {
		tree.entries = append(tree.entries, entryFromBytes(entryData[i], entryData[i+1]))
	}

	return tree, nil
}

func ListTree(args []string) {
	var nameOnly bool
	flagset := flag.NewFlagSet("ls-tree", flag.ExitOnError)
	flagset.BoolVar(&nameOnly, "name-only", false, "Only show names")
	err := flagset.Parse(args)
	die.If(err)

	if flagset.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "Usage: ls-tree [options] <object>")
		flagset.PrintDefaults()

		os.Exit(1)
	}

	id := flagset.Arg(0)
	obj, err := ReadBlobWithID(id)
	die.If(err)

	if obj.Type != TypeTree {
		fmt.Fprintf(os.Stderr, "%s is not a tree\n", id)
		os.Exit(1)
	}

	tree, err := TreeFromBlob(obj)
	die.If(err)

	if nameOnly {
		fmt.Print(tree.Names())
	} else {
		fmt.Print(tree)
	}
}

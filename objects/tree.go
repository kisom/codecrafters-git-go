package objects

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/die"
	"os"
	"sort"
	"strings"
)

type entryMode int

const (
	ModeRegular entryMode = (iota + 1) << 1
	ModeExecutable
	ModeSymbolic
	ModeDirectory
)

var modeToString = map[entryMode]string{
	ModeRegular:    "100644",
	ModeExecutable: "100755",
	ModeSymbolic:   "120000",
	ModeDirectory:  "40000",
}

var modeFromString = map[string]entryMode{
	"100644": ModeRegular,
	"100755": ModeExecutable,
	"120000": ModeSymbolic,
	"040000": ModeDirectory,
	"40000":  ModeDirectory,
}

type TreeEntry struct {
	Hash []byte
	Name string
	Mode entryMode
}

func (e *TreeEntry) Raw() []byte {
	// <mode> <name>\0<20_byte_sha>
	raw := []byte(fmt.Sprintf("%s %s\x00", modeToString[e.Mode], e.Name))
	raw = append(raw, e.Hash...)
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
	if len(b) != 20 {
		panic("bad hash length")
	}
	entry := &TreeEntry{}
	s := strings.Fields(string(a))
	entry.Mode = modeFromString[s[0]]
	entry.Name = s[1]
	entry.Hash = b
	return entry
}

type treeEntries []*TreeEntry

func (e treeEntries) Len() int      { return len(e) }
func (e treeEntries) Swap(i, j int) { e[i], e[j] = e[j], e[i] }
func (e treeEntries) Less(i, j int) bool {
	//if e[i].Mode != e[j].Mode {
	//	if e[i].Mode == ModeDirectory {
	//		return true
	//	} else if e[j].Mode == ModeDirectory {
	//		return false
	//	}
	//}
	return e[i].Name < e[j].Name
}

type Tree struct {
	entries treeEntries
}

func (tree *Tree) Add(entry *TreeEntry) {
	for _, e := range tree.entries {
		if e.Name == entry.Name {
			return
		}
	}

	tree.entries = append(tree.entries, entry)
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

func (tree *Tree) Write() error {
	blob := BlobFromBytes(tree.rawEntries())
	blob.Type = TypeTree

	if blob.HashString() != tree.HashString() {
		fmt.Fprintf(os.Stderr, "tree.Write: blob hash %s != tree hash %s", blob.HashString(), tree.HashString())
	}

	return blob.Write()
}

func scanEntries(contents []byte) []*TreeEntry {
	line := []byte{}
	entries := []*TreeEntry{}

	for i := 0; i < len(contents); i++ {
		if contents[i] != '\x00' {
			line = append(line, contents[i])
			continue
		}

		i++
		entry := entryFromBytes(line, contents[i:i+20])
		line = []byte{}
		i += 19
		entries = append(entries, entry)
	}

	if len(line) > 0 {
		panic("leftover data: " + string(line))
	}
	return entries
}

func TreeFromBlob(blob *Blob) (*Tree, error) {
	tree := &Tree{}
	tree.entries = scanEntries(blob.Contents)
	sort.Sort(tree.entries)

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

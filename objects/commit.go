package objects

import (
	"fmt"
	"git.wntrmute.dev/kyle/goutils/die"
	"github.com/kisom/codecrafters/git-go/paths"
	"os"
	"strings"
	"time"
)

const TypeCommit = "commit"

type Commit struct {
	Author    string
	Timestamp time.Time
	Tree      string
	Parents   []string
	Message   string
}

func (c *Commit) parents() string {
	buf := ""
	for _, parent := range c.Parents {
		buf += "parent " + parent + "\n"
	}

	return strings.TrimSpace(buf)
}

func (c *Commit) String() string {
	return fmt.Sprintf("tree %s\n%s\nauthor %s\ncommiter %s\n\n%s\n",
		c.Tree, c.parents(), AuthorLine(c.Author, c.Timestamp), AuthorLine(c.Author, c.Timestamp), c.Message)
}

func (c *Commit) blob() *Blob {
	return &Blob{
		Type:     TypeCommit,
		Contents: []byte(c.String()),
	}
}

func (c *Commit) Hash() []byte {
	return c.blob().Hash()
}

func (c *Commit) HashString() string {
	return c.blob().HashString()
}

func (c *Commit) Write() error {
	return c.blob().Write()
}

func NewCommitFromTree(tree, parent, message string) *Commit {
	commit := &Commit{
		Author:    DefaultAuthor(),
		Timestamp: time.Now(),
		Tree:      tree,
		Message:   message,
	}

	if parent != "" {
		commit.Parents = append(commit.Parents, parent)
	}

	return commit
}

func CommitTree(args []string) {
	var parent, message string
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: commit [tree] [-p parent] [-m message]\n")
		os.Exit(1)
	}

	treeID := args[0]
	args = args[1:]

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-p":
			parent = args[i+1]
		case "-m":
			message = args[i+1]
		default:
			fmt.Fprintf(os.Stderr, "Unknown argument '%s'\n", args[i])
			fmt.Fprintf(os.Stderr, "Usage: commit [tree] [-p parent] [-m message]\n")
			os.Exit(1)
		}
		i++
	}

	if message == "" {
		fmt.Fprintf(os.Stderr, "missing commit message\n")
		fmt.Fprintf(os.Stderr, "Usage: %s commit-tree <tree_sha> -m <message>\n", os.Args[0])
		os.Exit(1)
	}

	commit := NewCommitFromTree(treeID, parent, message)
	commit.Author = DefaultAuthor()

	err := commit.Write()
	die.If(err)

	err = paths.WriteRef("heads/master", commit.HashString())
	die.If(err)

	fmt.Println(commit.HashString())
}

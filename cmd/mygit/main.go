package main

import (
	"fmt"
	"git.wntrmute.dev/kyle/goutils/die"
	"github.com/kisom/codecrafters/git-go/git"
	"github.com/kisom/codecrafters/git-go/objects"
	"os"
	// Uncomment this block to pass the first stage!
	// "os"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	fmt.Fprintf(os.Stderr, "Program invocation: %#v\n", os.Args)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "cat-file":
		objects.CatFile(os.Args[2:])
	case "commit-tree":
		objects.CommitTree(os.Args[2:])
	case "hash-object":
		objects.Hash(os.Args[2:])
	case "ls-tree":
		if len(os.Args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
			os.Exit(1)
		}
		objects.ListTree(os.Args[2:])
	case "write-tree":
		hash, err := git.WriteTree()
		die.If(err)
		fmt.Println(hash)
	case "init":
		for _, dir := range []string{".git", ".git/objects", ".git/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/master\n")
		if err := os.WriteFile(".git/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized git directory")

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}

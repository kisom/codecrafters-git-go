package main

import (
	"flag"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/die"
	"git.wntrmute.dev/kyle/goutils/log"
	"os"

	"github.com/kisom/codecrafters/git-go/git"
	"github.com/kisom/codecrafters/git-go/objects"
)

// Usage: your_git.sh <command> <arg1> <arg2> ...
func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	// fmt.Println("Logs from your program will appear here!")

	opts := log.DefaultOptions("mygit", false)
	flag.StringVar(&opts.Level, "l", "info", "log level")
	flag.Parse()
	log.Setup(opts)

	args := flag.Args()

	if len(args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := args[0]; command {
	case "cat-file":
		objects.CatFile(args[1:])
	case "clone":
		git.Clone(args[1:])
	case "commit-tree":
		objects.CommitTree(args[1:])
	case "hash-object":
		objects.Hash(args[1:])
	case "ls-tree":
		if len(args) < 3 {
			fmt.Fprintf(os.Stderr, "usage: mygit <command> [<args>...]\n")
			os.Exit(1)
		}
		objects.ListTree(args[1:])
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

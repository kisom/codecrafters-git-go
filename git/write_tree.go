package git

import (
	"fmt"
	"github.com/kisom/codecrafters/git-go/objects"
	"github.com/kisom/codecrafters/git-go/paths"
	"github.com/pkg/errors"
	"io/fs"
	"os"
	"path/filepath"
)

func writeTree(path string) ([]byte, error) {
	dirents, err := os.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "writeTree")
	}

	tree := &objects.Tree{}

	for _, dirent := range dirents {
		entry := &objects.TreeEntry{}
		target := filepath.Join(path, dirent.Name())

		switch {
		case dirent.IsDir():
			switch dirent.Name() {
			case ".git", ".", "..":
				continue
			}

			hash, err := writeTree(target)
			if err != nil {
				return nil, err
			}

			entry.Hash = hash
			entry.Name = dirent.Name()
			entry.Mode = objects.ModeDirectory

		case dirent.Type().IsRegular():
			blob, err := objects.NewBlobFromFile(target)
			if err != nil {
				return nil, errors.Wrap(err, "writeTree")
			}
			err = blob.Write()
			if err != nil {
				return nil, errors.Wrap(err, "writeTree")
			}

			entry.Hash = blob.Hash()
			entry.Name = dirent.Name()
			entry.Mode = objects.ModeRegular
			if dirent.Type().Perm()&0111 != 0 {
				fmt.Fprintf(os.Stderr, "executable file has permissions %o\n", dirent.Type().Perm())
				entry.Mode = objects.ModeExecutable
			}

		case dirent.Type()&fs.ModeSymlink != 0:
			target, err := os.Readlink(filepath.Join(path, dirent.Name()))
			if err != nil {
				return nil, errors.Wrap(err, "writeTree")
			}

			blob, err := objects.NewBlobFromFile(filepath.Join(path, target))
			if err != nil {
				return nil, errors.Wrap(err, "writeTree")
			}
			err = blob.Write()
			if err != nil {
				return nil, errors.Wrap(err, "writeTree")
			}

			entry.Hash = blob.Hash()
			entry.Name = dirent.Name()
			entry.Mode = objects.ModeSymbolic
		default:
			panic("unhandled dirent case")
		}

		tree.Add(entry)
	}

	err = tree.Write()
	if err != nil {
		return nil, errors.Wrap(err, "writeTree")
	}

	return tree.Hash(), nil
}

func WriteTree() (string, error) {
	topLevel, err := paths.FindGitRoot()
	if err != nil {
		return "", errors.Wrap(err, "writing tree")
	}

	// Strip out the .git repository.
	hash, err := writeTree(topLevel)
	if err != nil {
		return "", errors.Wrap(err, "writing tree")
	}

	return fmt.Sprintf("%02x", hash), nil
}

package paths

import (
	"fmt"
	"git.wntrmute.dev/kyle/goutils/assert"
	"git.wntrmute.dev/kyle/goutils/fileutil"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

const ObjectIDLength = 40

func PathFromID(id string) (string, error) {
	if len(id) != ObjectIDLength {
		return "", fmt.Errorf("object ID `%s` want length %d, have %d", id, ObjectIDLength, len(id))
	}
	assert.Bool(len(id) == ObjectIDLength)

	subdir := id[:2]
	base := id[2:]

	parent, err := FindGitRoot()
	if err != nil {
		return "", errors.Wrap(err, "while searching for git root trying to build PathFromID("+id+")")
	}

	return filepath.Join(parent, ".git", "objects", subdir, base), nil
}

func PathForRef(ref string) (string, error) {
	parent, err := FindGitRoot()
	if err != nil {
		return "", errors.Wrap(err, "while searching for git root trying to build PathForRef("+ref+")")
	}

	return filepath.Join(parent, ".git", "refs", ref), nil
}

func WriteRef(ref string, id string) error {
	refpath, err := PathForRef(ref)
	if err != nil {
		return errors.Wrap(err, "while writing ref")
	}

	parent := filepath.Dir(refpath)
	if !fileutil.DirectoryDoesExist(parent) {
		err = os.MkdirAll(parent, 0755)
		if err != nil {
			return errors.Wrap(err, "while writing ref")
		}
	}

	file, err := os.Create(refpath)
	if err != nil {
		return errors.Wrap(err, "while writing ref")
	}
	defer file.Close()

	fmt.Fprintf(file, "%s\n", id)
	return nil
}

func FindGitRoot() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", errors.Wrap(err, "unable to determine current working directory")
	}
	startingDirectory := cwd
	defer func(dir string) {
		err := os.Chdir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to return starting directory %s\n", startingDirectory)
		}
	}(startingDirectory)

	for {
		if fileutil.DirectoryDoesExist(filepath.Join(cwd, ".git")) {
			return cwd, nil
		}

		if cwd == "/" {
			return "", errors.New("repository not found; stopped searching at root directory")
		}

		if err = os.Chdir(".."); err != nil {
			return "", errors.Wrap(err, "unable to change directory to parent")
		}

		cwd, err = os.Getwd()
		if err != nil {
			return "", errors.Wrap(err, "unable to determine current working directory")
		}
	}
}

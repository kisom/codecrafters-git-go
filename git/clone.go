package git

import (
	"flag"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/fileutil"
	"git.wntrmute.dev/kyle/goutils/log"
	"github.com/kisom/codecrafters/git-go/pack"
	"github.com/pkg/errors"
	"net/url"
	"os"
	"path/filepath"
)

func defaultRepoDir(repo string) (string, error) {
	repoURL, err := url.Parse(repo)
	if err != nil {
		return "", errors.Wrap(err, "parsing repo URL")
	}

	log.Debugf("repo URL: %#v\n", repoURL)

	dirName := filepath.Base(repoURL.Path)
	if dirName == "" {
		return repoURL.Hostname(), nil
	}

	return dirName, nil
}

func removeDirectoryIfEmpty(dir string) {
	log.Debugf("removing directory %s\n", dir)
	if !fileutil.FileDoesExist(filepath.Join(dir, ".git", "HEAD")) {
		err := os.RemoveAll(dir)
		if err != nil {
			log.Errf("repository is empty but failed to remove directory %s: %v", dir, err)
		}
	}
}

func clonePreflight(repo string, dirName string) error {
	fmt.Printf("cloning repository %s into %s\n", repo, dirName)

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error creating repository directory:", err)
		return err
	}

	log.Debugf("repository directory created: %s\n", dirName)

	return nil
}

func cloneReferences(dirName string, ra *pack.ReferenceAdvertisement) error {
	return errors.New("not implemented")
}

func clone(repo string, dirName string) error {
	err := clonePreflight(repo, dirName)
	if err != nil {
		return err
	}
	defer removeDirectoryIfEmpty(dirName)

	advertisement, err := pack.FetchReferenceAdvertisement(repo)
	if err != nil {
		return err
	}

	log.Infof("read %d reference in advertisement", len(advertisement.References))

	return nil
}

func Clone(args []string) {
	flags := flag.NewFlagSet("clone", flag.ExitOnError)
	err := flags.Parse(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "no repo provided.")
		os.Exit(1)
	}

	repo := flags.Arg(0)
	dirName, err := defaultRepoDir(repo)
	if len(args) > 1 {
		dirName = flags.Arg(1)
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "invalid repo provided: %s\n", err)
		os.Exit(1)
	}

	if fileutil.DirectoryDoesExist(dirName) {
		fmt.Fprintln(os.Stderr, "repository already exists in", dirName)
		os.Exit(1)
	}

	err = clone(repo, dirName)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error cloning repository:", err)
		os.Exit(1)
	}
}

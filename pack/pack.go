package pack

import (
	"bytes"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/log"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"path"
)

const (
	serviceGitUploadPackAdvertisement = "git-upload-pack-advertisement"
	serviceGitUploadPack              = "git-upload-pack"
)

func serviceContentType(service string) string {
	return fmt.Sprintf("application/x-%s", service)
}

func normalizeRepoURL(repo string) (*url.URL, error) {
	repoURL, err := url.Parse(repo)
	if err != nil {
		return nil, errors.Wrap(err, "parsing repo URL")
	}

	if repoURL.Scheme != "http" && repoURL.Scheme != "https" {
		repoURL.Scheme = "https"
	}

	return repoURL, nil
}

func serviceURL(repo, service string) (string, error) {
	repoURL, err := normalizeRepoURL(repo)
	if err != nil {
		return "", errors.Wrap(err, "normalizing repo URL")
	}

	repoURL.Path = path.Join(repoURL.Path, "info", "refs")
	repoURL.RawQuery = "service=" + service

	log.Debugf("repo service URL: %s", repoURL.String())

	return repoURL.String(), nil
}

func FetchReferenceAdvertisement(repo string) (*ReferenceAdvertisement, error) {
	repoURL, err := serviceURL(repo, serviceGitUploadPackAdvertisement)
	if err != nil {
		return nil, errors.Wrap(err, "normalizing repo URL")
	}

	resp, err := http.Get(repoURL)
	if err != nil {
		return nil, errors.Wrap(err, "fetching repo packfile")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, fmt.Errorf("fetching repo packfile returned http status code %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != serviceContentType(serviceGitUploadPack) {
		return nil, fmt.Errorf("unsupported content type %q (expect %s)",
			resp.Header.Get("Content-Type"), serviceContentType(serviceGitUploadPack))
	}

	packfile, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading packfile")
	}

	if !referenceAdvertisementMagic.Match(packfile[:5]) {
		return nil, fmt.Errorf("packfile contains invalid magic %x", packfile[:5])
	}

	ra := &ReferenceAdvertisement{}
	err = ra.UnmarshalReader(bytes.NewBuffer(packfile))
	if err != nil {
		return nil, errors.Wrap(err, "parsing repo packfile body")
	}
	return ra, nil
}

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

const contentTypeServicePack = "application/x-git-upload-pack-advertisement"

func FetchReferenceAdvertisement(repo string) (*ReferenceAdvertisement, error) {
	repoURL, err := url.Parse(repo)
	if err != nil {
		return nil, errors.Wrap(err, "parsing repo URL")
	}

	if repoURL.Scheme != "http" && repoURL.Scheme != "https" {
		repoURL.Scheme = "https"
	}

	repoURL.Path = path.Join(repoURL.Path, "info", "refs")
	repoURL.RawQuery = "service=git-upload-pack"

	log.Debugf("Fetching packfile from %s", repoURL.String())

	resp, err := http.Get(repoURL.String())
	if err != nil {
		return nil, errors.Wrap(err, "fetching repo packfile")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		return nil, fmt.Errorf("fetching repo packfile returned http status code %d", resp.StatusCode)
	}

	if resp.Header.Get("Content-Type") != contentTypeServicePack {
		return nil, fmt.Errorf("unsupported content type %q (expect %s)",
			resp.Header.Get("Content-Type"), contentTypeServicePack)
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

package pack

import (
	"os"
	"testing"
)

func TestFetchPackfile(t *testing.T) {
	if os.Getenv("MYGIT_TEST_NET_ENABLED") != "Y" {
		t.Skip()
	}

	repo := "github.com/kisom/codecrafters-git-go"
	_, err := FetchReferenceAdvertisement(repo)
	if err != nil {
		t.Fatalf("FetchReferenceAdvertisement failed: %s", err)
	}

}

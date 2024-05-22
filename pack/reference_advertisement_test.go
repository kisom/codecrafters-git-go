package pack

import (
	"bytes"
	"os"
	"testing"
)

func TestUnmarshalReferenceAdvertisement(t *testing.T) {
	file, err := os.Open("testdata/packfile")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	ra := &ReferenceAdvertisement{}
	err = ra.UnmarshalReader(file)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUnmarshalReference(t *testing.T) {
	testReference := bytes.NewBufferString("0155378dced80975bb1e0ccf3aac86be65f058683e6c HEAD\x00multi_ack thin-pack side-band side-band-64k ofs-delta shallow deepen-since deepen-not deepen-relative no-progress include-tag multi_ack_detailed allow-tip-sha1-in-want allow-reachable-sha1-in-want no-done symref=HEAD:refs/heads/master filter object-format=sha1 agent=git/github-f133c3a1d7e6\n003f378dced80975bb1e0ccf3aac86be65f058683e6c refs/heads/master\n0000")

	ra := &ReferenceAdvertisement{}
	err := ra.readReferences(testReference)
	if err != nil {
		t.Fatal(err)
	}
}

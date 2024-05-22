package pack

import (
	"bytes"
	"os"
	"testing"
)

const testPackFile = "testdata/packfile"

func TestReadPktLine(t *testing.T) {
	expectedFirstLine := []byte("# service=git-upload-pack")
	file, err := os.Open(testPackFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	pktLine, err := readPktLine(file)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pktLine, expectedFirstLine) {
		t.Fatalf("invalid first line\n\texpected: '%s'\n\t    have: '%s'\n",
			expectedFirstLine, pktLine)
	}
}

func TestWritePktLineString(t *testing.T) {
	expectedPktLine := []byte("001e# service=git-upload-pack\n")
	payload, err := readPktLine(bytes.NewBuffer(expectedPktLine))
	if err != nil {
		t.Fatal(err)
	}

	pktLine, err := writePacketLineString(string(payload))
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(pktLine, expectedPktLine) {
		t.Fatalf("invalid first line\n\texpected: '%s'\n\t    have: '%s'\n", expectedPktLine, pktLine)
	}
}

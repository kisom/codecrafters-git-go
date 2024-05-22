package pack

import (
	"bytes"
	"fmt"
	"io"
	"regexp"
	"strconv"

	"github.com/pkg/errors"
)

// note: the git docs call this type "pkt-line," not "packet-line." The code
// follows this convention.

var flushPkt = regexp.MustCompile(`0000$`)

func readPktLine(r io.Reader) ([]byte, error) {
	length := make([]byte, 4)
	_, err := r.Read(length)
	if err != nil {
		return nil, err
	}

	lineLength, err := strconv.ParseInt(string(length), 16, 16)
	if err != nil {
		return nil, errors.Wrap(err, "parse line length")
	}

	if lineLength < 5 {
		if lineLength == 0 {
			return nil, nil // nil signifies a flush-pkt
		}
		return nil, errors.New("empty pkt-line sent")
	}
	lineLength -= 4

	buf := make([]byte, lineLength)
	n, err := r.Read(buf)
	if err != nil {
		return nil, errors.Wrap(err, "read pkt-line")
	}

	if n != int(lineLength) {
		return nil, fmt.Errorf("read pkt-line too small; have %d, want %d", n, lineLength)
	}

	return bytes.TrimSpace(buf), nil
}

func writePacketLineString(line string) ([]byte, error) {
	length := len(line)
	length += 5

	return []byte(fmt.Sprintf("%04x%s\n", length, line)), nil
}

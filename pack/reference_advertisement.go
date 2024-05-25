package pack

import (
	"bytes"
	"fmt"
	"git.wntrmute.dev/kyle/goutils/log"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"strings"
)

var referenceAdvertisementMagic = regexp.MustCompile(`^[0-9a-f]{4}#$`)

type Reference struct {
	ID           string
	Name         string
	Capabilities []string
}

func (ref *Reference) Want() []byte {
	return writePacketLineString("want " + ref.ID)
}

func (ref *Reference) Have() []byte {
	return writePacketLineString("have " + ref.ID)
}

type ReferenceAdvertisement struct {
	References []*Reference
}

func (ra *ReferenceAdvertisement) UnmarshalReader(r io.Reader) error {
	advertisementLine, err := readPktLine(r)
	if err != nil {
		return errors.Wrap(err, "reading advertisement line")
	}

	if !bytes.Equal(advertisementLine, []byte("# service=git-upload-pack")) {
		return fmt.Errorf("advertisement line is invalid")
	}

	log.Debugln("advertisement line is good")

	// should be a flush-pkt
	_, err = readPktLine(r)
	if err != nil {
		return errors.Wrap(err, "reading advertisement line")
	}

	err = ra.readReferences(r)
	if err != nil {
		return errors.Wrap(err, "reading references")
	}

	return nil
}

func (ra *ReferenceAdvertisement) Want(w io.Writer) error {
	for _, ref := range ra.References {
		_, err := w.Write(ref.Want())
		if err != nil {
			return errors.Wrap(err, "writing reference advertisement")
		}
	}

	_, err := w.Write(flushPkt)
	return err
}

func (ra *ReferenceAdvertisement) readReferences(r io.Reader) error {
	for {
		line, err := readPktLine(r)
		if err != nil {
			return errors.WithStack(errors.Wrap(err, "reading references line"))
		}

		if len(line) == 0 {
			break
		}

		ref, err := parseReferenceLine(line)
		if err != nil {
			return errors.Wrap(err, "parsing references line")
		}

		ra.References = append(ra.References, ref)
	}

	return nil
}

func parseReferenceLine(line []byte) (*Reference, error) {
	ref := &Reference{}
	fields := bytes.Split(line, []byte{0})
	if len(fields) > 2 {
		return nil, fmt.Errorf("invalid reference line; too many (%d) fields", len(fields))
	}

	refDescription := strings.Fields(string(fields[0]))
	ref.ID = strings.TrimSpace(refDescription[0])
	ref.Name = strings.TrimSpace(refDescription[1])

	if len(fields) > 1 {
		ref.Capabilities = strings.Fields(string(bytes.TrimSpace(fields[1])))
	}

	return ref, nil
}

package wts

import (
	"encoding/binary"
	"io"

	"github.com/cockroachdb/errors"
)

type FrameFlags byte

const (
	FlagCompressed FrameFlags = 1 << iota
	FlagEncrypted
	FlagFragment
	FlagNone FrameFlags = 0
)

type Header struct {
	Magic   uint16
	Version byte
	Type    MessageType
	Flags   FrameFlags
	Length  uint32
}

func (h *Header) Encode() []byte {

	buf := make([]byte, HeaderSize)

	binary.BigEndian.PutUint16(
		buf[0:2],
		h.Magic,
	)

	buf[2] = h.Version

	buf[3] = byte(h.Type)

	buf[4] = byte(h.Flags)

	binary.BigEndian.PutUint32(
		buf[5:9],
		h.Length,
	)

	return buf
}

func DecodeHeader(buf []byte) (*Header, error) {

	if len(buf) < HeaderSize {
		return nil, errors.New("invalid header size")
	}

	h := &Header{

		Magic: binary.BigEndian.Uint16(buf[0:2]),

		Version: buf[2],

		Type: MessageType(buf[3]),

		Flags: FrameFlags(buf[4]),

		Length: binary.BigEndian.Uint32(buf[5:9]),
	}

	if err := validateHeader(h); err != nil {
		return nil, err
	}

	return h, nil
}

func Write(w io.Writer, t MessageType, payload []byte) error {

	header := Header{
		Magic:   Magic,
		Version: Version,
		Type:    t,
		Flags:   FlagNone,
		Length:  uint32(len(payload)),
	}

	if _, err := w.Write(header.Encode()); err != nil {
		return err
	}

	_, err := w.Write(payload)
	return err
}

func Read(r io.Reader) (*Header, []byte, error) {

	headerBuf := make([]byte, HeaderSize)

	if _, err := io.ReadFull(r, headerBuf); err != nil {
		return nil, nil, err
	}

	header, err := DecodeHeader(headerBuf)
	if err != nil {
		return nil, nil, err
	}

	if header.Length == 0 {
		return header, nil, nil
	}

	payload := make([]byte, header.Length)

	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, nil, err
	}
	return header, payload, nil
}

func validateHeader(h *Header) error {
	if h.Magic != Magic {
		return errors.New("invalid magic")
	}

	if h.Version != Version {
		return errors.New("unsupported version")
	}

	// TODO: supported flags
	// h.Flags
	return nil
}

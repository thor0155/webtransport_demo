package codectool

import (
	"github.com/cockroachdb/errors"
	"github.com/vmihailenco/msgpack/v5"
)

func Encode(
	msg any,
) ([]byte, error) {
	if b, err := msgpack.Marshal(msg); err != nil {
		return nil, errors.WithStack(err)
	} else {
		return b, nil
	}
}

func Decode(
	data []byte,
	ptr any,
) error {
	return errors.WithStack(msgpack.Unmarshal(data, ptr))
}

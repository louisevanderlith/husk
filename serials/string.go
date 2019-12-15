package serials

import (
	"io"
)

type StringSerial struct {
}

func (s StringSerial) Encode(obj interface{}) ([]byte, error) {
	return []byte(obj.(string)), nil
}

func (s StringSerial) Decode(r io.Reader, obj interface{}) error {
	var b []byte
	r.Read(b)

	obj = string(b)

	return nil
}

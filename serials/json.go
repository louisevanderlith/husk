package serials

import (
	"encoding/json"
	"io"
)

type JsonSerial struct {
}

func (s JsonSerial) Encode(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}

func (s JsonSerial) Decode(r io.Reader, obj interface{}) error {
	return json.NewDecoder(r).Decode(obj)
}

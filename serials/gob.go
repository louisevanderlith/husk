package serials

import (
	"bytes"
	"encoding/gob"
	"io"
)

type GobSerial struct {
}

func (s GobSerial) Encode(obj interface{}) ([]byte, error) {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(&obj)

	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (s GobSerial) Decode(r io.Reader, obj interface{}) error {
	d := gob.NewDecoder(r)
	err := d.Decode(obj)

	if err != nil {
		return err
	}

	return nil
}

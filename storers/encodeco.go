package storers

import (
	"encoding/gob"
	"io"
)

type NewDecoder func(r io.Reader) Decoder
type NewEncoder func(w io.Writer) Encoder

type Encoder interface {
	Encode(e interface{}) error
}

type Decoder interface {
	Decode(e interface{}) error
}

func GobEncoder(w io.Writer) Encoder {
	return gob.NewEncoder(w)
}

func GobDecoder(r io.Reader) Decoder {
	return gob.NewDecoder(r)
}

package storage

import (
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
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

func JSONEncoder(w io.Writer) Encoder {
	return json.NewEncoder(w)
}

func JSONDecoder(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func XMLEncoder(w io.Writer) Encoder {
	return xml.NewEncoder(w)
}

func XMLDecoder(r io.Reader) Decoder {
	return xml.NewDecoder(r)
}

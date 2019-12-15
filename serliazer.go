package husk

import "io"

type Serializer interface {
	Encode(obj interface{}) ([]byte, error)
	Decode(r io.Reader, obj interface{}) error
}

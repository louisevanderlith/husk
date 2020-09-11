package hsk

import (
	"github.com/louisevanderlith/husk/validation"
)

//Storage stores meta with some knowledge of data, a wrapper for a reader/writer
type Storage interface {
	Read(p Point, data chan<- Record)
	Write(obj Record, p chan<- Point)
	Close() error
	Name() string
	ZeroValue() validation.Dataer
}

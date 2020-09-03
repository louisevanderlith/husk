package storers

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/validation"
)

//Storage stores meta with knowledge of data, a wrapper for a reader/writer
type Storage interface {
	Read(p hsk.Point, data chan<- validation.Dataer)
	Write(obj validation.Dataer, p chan<- hsk.Point)
	Close() error
}

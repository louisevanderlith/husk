package storers

import (
	"github.com/louisevanderlith/husk/hsk"
)

//Storage stores meta with knowledge of data, a wrapper for a reader/writer
type Storage interface {
	Read(p hsk.Point, data chan<- hsk.Dataer) error
	Write(obj hsk.Dataer) (hsk.Point, error)
	Close() error
}

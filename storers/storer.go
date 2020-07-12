package storers

import (
	"github.com/louisevanderlith/husk/hsk"
)

//Storer stores meta with knowledge of data
type Storer interface {
	Read(p *hsk.Point, data chan<- hsk.Dataer) error
	Write(obj hsk.Dataer) (*hsk.Point, error)
	Close() error
}

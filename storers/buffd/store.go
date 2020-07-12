package buffd

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
)

func newStore() storers.Storer {
	return buffdStore{}
}

type buffdStore struct {
}

func (b buffdStore) Read(p *hsk.Point, res chan<- hsk.Dataer) error {
	panic("implement me")
}

func (b buffdStore) Write(obj hsk.Dataer) (*hsk.Point, error) {
	panic("implement me")
}

func (b buffdStore) Close() error {
	panic("implement me")
}

package hybrid

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
)

func newStore() storers.Storage {
	return &hybridStore{}
}

type hybridStore struct {
}

func (h hybridStore) Read(p hsk.Point, data chan<- hsk.Dataer) error {
	panic("implement me")
}

func (h hybridStore) Write(obj hsk.Dataer) (hsk.Point, error) {
	panic("implement me")
}

func (h hybridStore) Close() error {
	panic("implement me")
}

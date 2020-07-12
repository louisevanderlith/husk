package chip

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
)

func newStore() storers.Storer {
	return &chipStore{}
}

type chipStore struct {
	records []hsk.Dataer
}

func (c *chipStore) Read(p *hsk.Point, res chan<- hsk.Dataer) error {
	go func() {
		res <- c.records[p.Offset]
	}()

	return nil
}

func (c *chipStore) Write(obj hsk.Dataer) (*hsk.Point, error) {
	c.records = append(c.records, obj)

	ln := int64(len(c.records))
	return hsk.NewPoint(ln-1, ln), nil
}

func (c chipStore) Close() error {
	panic("implement me")
}

package sample

import (
	"github.com/louisevanderlith/husk/collections"
	"github.com/louisevanderlith/husk/db"
	"github.com/louisevanderlith/husk/op"
	"github.com/louisevanderlith/husk/storers"
	"github.com/louisevanderlith/husk/storers/chip"
)

type SampleContext interface {
	db.Ctxer
	FindJournals(page, size int) (collections.Page, error)
	FindJournalsByPublisher(page, size int, name string) (collections.Page, error)
	HasJournals() bool
	CountJournals() (int64, error)
}

type context struct {
	Journals storers.Table
}

//Returns a new Journal Database
func NewContext() SampleContext {
	result := context{
		Journals: chip.NewTable(Journal{}),
	}

	result.seed()
	return result
}

func (ctx context) seed() {
	err := ctx.Journals.Seed("journals.seed.json")

	if err != nil {
		panic(err)
	}
}

func (ctx context) Shutdown() error {
	return nil
}

func (ctx context) HasJournals() bool {
	return ctx.Journals.Exists(op.Everything())
}

func (ctx context) CountJournals() (int64, error) {
	count := int64(0)
	err := ctx.Journals.Map(&count, op.RowCount())

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (ctx context) FindJournals(page, size int) (collections.Page, error) {
	return ctx.Journals.Find(page, size, op.Everything())
}

func (ctx context) FindJournalsByPublisher(page, size int, name string) (collections.Page, error) {
	return ctx.Journals.Find(page, size, ByPublisher(name))
}

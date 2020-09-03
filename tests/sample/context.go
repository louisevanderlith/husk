package sample

import (
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/op"
)

type JournalContext interface {
	FindJournals(page, size int) (hsk.Page, error)
	FindJournalsByPublisher(page, size int, name string) (hsk.Page, error)
	HasJournals() bool
	CountJournals() (int64, error)
}

type context struct {
	Journals hsk.Table
}

//Returns a new Journal Database
func NewContext() JournalContext {
	result := context{
		Journals: husk.NewTable(Journal{}),
	}

	result.Seed()
	return result
}

func (ctx context) Seed() {
	err := ctx.Journals.Seed("journals.seed.json")

	if err != nil {
		panic(err)
	}
}

func (ctx context) Save() error {
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

func (ctx context) FindJournals(page, size int) (hsk.Page, error) {
	return ctx.Journals.Find(page, size, op.Everything())
}

func (ctx context) FindJournalsByPublisher(page, size int, name string) (hsk.Page, error) {
	return ctx.Journals.Find(page, size, ByPublisher(name))
}

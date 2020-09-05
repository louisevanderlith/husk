package sample

import (
	"encoding/json"
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/collections"
	"github.com/louisevanderlith/husk/op"
	"github.com/louisevanderlith/husk/records"
	"os"
	"reflect"
)

type JournalContext interface {
	FindJournals(page, size int) (records.Page, error)
	FindJournalsByPublisher(page, size int, name string) (records.Page, error)
	HasJournals() bool
	CountJournals() (int64, error)
}

type context struct {
	Journals husk.Table
}

//Returns a new Journal Database
func NewContext() JournalContext {
	result := context{
		Journals: husk.NewTable(Journal{}),
	}

	err := result.Seed()

	if err != nil {
		panic(err)
	}

	return result
}

func (ctx context) Seed() error {
	journals, err := ReadJournals()

	if err != nil {
		return err
	}

	return ctx.Journals.Seed(journals)
}

func ReadJournals() (collections.Enumerable, error) {
	f, err := os.Open("journals.seed.json")

	if err != nil {
		return nil, err
	}

	var journals []Journal
	dec := json.NewDecoder(f)
	err = dec.Decode(&journals)

	if err != nil {
		return nil, err
	}

	return collections.ReadOnlyList(reflect.ValueOf(journals)), nil
}

func (ctx context) Save() error {
	return ctx.Journals.Save()
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

func (ctx context) FindJournals(page, size int) (records.Page, error) {
	return ctx.Journals.Find(page, size, op.Everything())
}

func (ctx context) FindJournalsByPublisher(page, size int, name string) (records.Page, error) {
	return ctx.Journals.Find(page, size, ByPublisher(name))
}

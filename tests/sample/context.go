package sample

import (
	"github.com/louisevanderlith/husk"
)

type Context struct {
	Journals husk.Tabler
}

//Returns a new Journal Database
func NewContext() Context {
	result := Context{}

	result.Journals = husk.NewTable(Journal{})

	return result
}

func (ctx Context) Seed() {
	err := ctx.Journals.Seed("journals.seed.json")

	if err != nil {
		panic(err)
	}

	err = ctx.Journals.Save()

	if err != nil {
		panic(err)
	}
}

func (ctx Context) Save() error {
	return ctx.Journals.Save()
}

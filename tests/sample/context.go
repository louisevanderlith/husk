package sample

import (
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/serials"
)

type Context struct {
	People husk.Tabler
	//Names  husk.Tabler
}

func NewContext() Context {
	result := Context{}

	result.People = husk.NewTable(Person{}, &serials.GobSerial{})
	//result.Names = husk.NewTable(Name{}, &serials.StringSerial{})

	return result
}

func (ctx Context) Seed() {
	err := ctx.People.Seed("people.seed.json")

	if err != nil {
		panic(err)
	}

	ctx.People.Save()
}

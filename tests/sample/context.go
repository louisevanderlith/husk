package sample

import "github.com/louisevanderlith/husk"

type Context struct {
	People husk.Tabler
}

func NewContext() Context {
	result := Context{}

	result.People = husk.NewTable(Person{})

	return result
}

func (ctx Context) Seed() {
	err := ctx.People.Seed("people.seed.json")

	if err != nil {
		panic(err)
	}

	ctx.People.Save()
}

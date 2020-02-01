package sample

import (
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/serials"
)

type Context struct {
	People husk.Tabler
	Users  husk.Tabler
}

func NewContext() Context {
	result := Context{}

	result.People = husk.NewTable(Person{}, serials.GobSerial{})
	result.Users = husk.NewTable(User{}, serials.GobSerial{})

	return result
}

func (ctx Context) Seed() {
	err := ctx.Users.Seed("users.seed.json")

	if err != nil {
		panic(err)
	}

	ctx.Users.Save()

	err = ctx.People.Seed("people.seed.json")

	if err != nil {
		panic(err)
	}

	ctx.People.Save()
}

func (ctx Context) Save() error {
	return ctx.People.Save()
}

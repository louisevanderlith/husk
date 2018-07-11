package sample

import "github.com/louisevanderlith/husk"

type Context struct {
	People   husk.Tabler
	Accounts husk.Tabler
}

func NewContext() Context {
	result := Context{}
	result.People = husk.NewTable(new(Person))
	result.Accounts = husk.NewTable(new(Account))

	return result
}

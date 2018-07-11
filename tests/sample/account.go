package sample

import "github.com/louisevanderlith/husk"

type Account struct {
	Number string
	Person *Person
}

func (o Account) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}

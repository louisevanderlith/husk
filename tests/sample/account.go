package sample

import "github.com/louisevanderlith/husk"

type Account struct {
	Number  string
	Balance float32
}

func (o Account) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}

package sample

import "github.com/louisevanderlith/husk"

type Person struct {
	Name     string `hsk:"size(50)"`
	Age      int
	Accounts []*Account
}

func (o Person) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}

package sample

import "github.com/louisevanderlith/husk"

type Role struct {
	ClientID string `hsk:"size(50)"`
	Scope    string `hsk:"size(50)"`
	Claim    string `hsk:"size(50)"`
	Value    string `hsk:"size(100)"`
}

func (o Role) Valid() error {
	return husk.ValidateStruct(&o)
}

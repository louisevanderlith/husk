package sample

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/validation"
)

type Event struct {
	Type     string
	Relation hsk.Key
}

func (e Event) Valid() error {
	return validation.Struct(e)
}

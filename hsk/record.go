package hsk

import (
	"github.com/louisevanderlith/husk/validation"
)

//Record is what defines a record, and what it can do
type Record interface {
	GetKey() Key
	GetValue() validation.Dataer
}

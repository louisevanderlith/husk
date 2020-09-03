package op

import (
	"github.com/louisevanderlith/husk/validation"
)

// Folder used to fold records, which can be used to perform Aggregate functions
type Folder interface {
	Fold(obj validation.Dataer) bool
}

type folder func(obj validation.Dataer) bool

//Fold is the function which performs aggregate functions on data.
func (f folder) Fold(obj validation.Dataer) bool {
	return f(obj)
}

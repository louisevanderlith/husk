package op

import (
	"github.com/louisevanderlith/husk/hsk"
)

// Folder used to fold records, which can be used to perform Aggregate functions
type Folder interface {
	Fold(obj hsk.Dataer) bool
}

type folder func(obj hsk.Dataer) bool

//Fold is the function which performs aggregate functions on data.
func (f folder) Fold(obj hsk.Dataer) bool {
	return f(obj)
}

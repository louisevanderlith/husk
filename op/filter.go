package op

import (
	"github.com/louisevanderlith/husk/hsk"
	"reflect"
)

// Filterer used to filter records while searching.
type Filterer interface {
	// Enables a single central function to cast from husk.Dataer to *<Person>
	Filter(obj hsk.Dataer) bool
}

type filter func(obj hsk.Dataer) bool

//Filter is the function which casts objects before sending them to the filter.
func (f filter) Filter(obj hsk.Dataer) bool {
	return f(obj)
}

// Everything, returns 'true' on all rows
func Everything() filter {
	return func(obj hsk.Dataer) bool {
		return true
	}
}

// ByFields returns objects that have properties that match with the given object
func ByFields(param hsk.Dataer) filter {
	parmFields := hsk.GetFields(param)
	return func(obj hsk.Dataer) bool {
		objFields := hsk.GetFields(obj)
		for k, v := range parmFields {
			if !reflect.ValueOf(v).IsZero() && !reflect.DeepEqual(objFields[k], v) {
				return false
			}
		}

		return true
	}
}

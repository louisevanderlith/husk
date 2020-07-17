package op

import (
	"github.com/louisevanderlith/husk/hsk"
	"reflect"
	"time"
)

// Filter used to filter records while searching.
type Filter interface {
	// Enables a single central function to cast from husk.Dataer to <Entity>
	Filter(obj hsk.Record) bool
}

type FilterFunc func(obj hsk.Record) bool

//Filter is the function which casts objects before sending them to the filter.
func (f FilterFunc) Filter(obj hsk.Record) bool {
	return f(obj)
}

// Everything, returns 'true' on all rows, inactive included
func Everything() FilterFunc {
	return func(obj hsk.Record) bool {
		return true
	}
}

// Everything, returns 'true' on all rows, created between the specified dates.
func EverythingBetween(start time.Time, end time.Time) FilterFunc {
	return func(obj hsk.Record) bool {
		stamp := obj.GetKey().GetTimestamp()

		return (stamp.After(start) || stamp.Equal(start)) && (stamp.Before(end) || stamp.Equal(end))
	}
}

// ByFields returns objects that have properties that match with the given object
func ByFields(param hsk.Dataer) FilterFunc {
	parmFields := hsk.GetFields(param)
	return func(obj hsk.Record) bool {
		objFields := hsk.GetFields(obj.Data())
		for k, v := range parmFields {
			if !reflect.ValueOf(v).IsZero() && !reflect.DeepEqual(objFields[k], v) {
				return false
			}
		}

		return true
	}
}

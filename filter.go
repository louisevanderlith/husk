package husk

import (
	"errors"
	"reflect"
)

type Filter func(obj Dataer) bool

func MakeFilter(customFilter interface{}) (Filter, error) {
	tpe := reflect.TypeOf(customFilter)

	if tpe.Kind() != reflect.Func {
		return nil, errors.New("customFilter is not a function.")
	}

	return func(obj Dataer) bool {
		val := reflect.ValueOf(customFilter)
		var params []reflect.Value
		params = append(params, reflect.ValueOf(obj))

		result := val.Call(params)

		return result[0].Bool()
	}, nil
}

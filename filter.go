package husk

type Filter func(obj Dataer) bool

// Ditto transforms your 'custom' filter into a husk.Filter
/*func Ditto(customFilter interface{}) Filter {
	val := transform(customFilter)

	return func(obj Dataer) bool {
		params := []reflect.Value{reflect.ValueOf(obj)}

		result := val.Call(params)

		return result[0].Bool()
	}
}

func transform(customFilter interface{}) reflect.Value {
	tpe := reflect.TypeOf(customFilter)

	if tpe.Kind() != reflect.Func {
		panic("customFilter is not a function.")
	}

	return reflect.ValueOf(customFilter)
}
*/

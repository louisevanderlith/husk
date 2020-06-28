package husk

import (
	"reflect"
)

// Ctxer is a special interface which can be used when building extensions
type Ctxer interface {
	//Save calls the save method for every registered table
	Save() error
}

// TableNames returns the names of all the tables found within the given context
func TableNames(ctx Ctxer) []string {
	var result []string

	val := reflect.ValueOf(ctx)
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		typeField := valType.Field(i)
		result = append(result, typeField.Name)
	}

	return result
}

// TableLayouts returns the names and table found within the given context
func TableLayouts(ctx Ctxer) map[string]Tabler {
	result := make(map[string]Tabler)

	val := reflect.ValueOf(ctx)
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := valType.Field(i)

		result[typeField.Name] = valueField.Interface().(Tabler)
	}

	return result
}

package validation

import (
	"errors"
	"reflect"
	"strings"
)

//Dataer is the primary interface that any "model" should implement
//"Models" are data objects used to store and structure records in tables.
type Dataer interface {
	Valid() error
}

// GetFields returns
func GetFields(t Dataer) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(t)
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := valType.Field(i)

		result[typeField.Name] = valueField.Interface()
	}

	return result
}

// ValidateStruct will read 'hsk' tags on properties, to validate their values
// Properties without the 'hsk' tag, will be considered 'Required'
func ValidateStruct(obj interface{}) error {
	var issues []string

	val := reflect.ValueOf(obj)
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := valType.Field(i)
		tag := typeField.Tag.Get("hsk")

		kind := valueField.Kind()
		validator := GetTypeValidator(kind)

		meta := GetMeta(tag, kind)
		meta.PropName = typeField.Name
		value := valueField.Interface()

		isValid, problems := validator.Valid(value, meta)

		if !isValid {
			issues = append(issues, problems...)
		}
	}

	if len(issues) != 0 {
		return errors.New(strings.Join(issues, "\r\n"))
	}

	return nil
}

package husk

import (
	"errors"
	"reflect"
	"strings"

	"github.com/louisevanderlith/husk/validation"
)

//Dataer is the primary interface that any "model" should implement
//"Models" are data objects used to store and structure records in tabes.
type Dataer interface {
	Valid() (bool, error)
}

// ValidateStruct will read 'hsk' tags on properties, to validate their values
// Properties without the 'hsk' tag, will be considered 'Required'
func ValidateStruct(obj interface{}) (bool, error) {
	var issues []string

	val := reflect.ValueOf(obj).Elem()
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := valType.Field(i)
		tag := typeField.Tag.Get("hsk")

		kind := valueField.Kind()
		validator := validation.GetTypeValidator(kind)

		meta := validation.GetMeta(tag, kind)
		meta.PropName = typeField.Name
		value := valueField.Interface()

		isValid, problems := validator.Valid(value, meta)

		if !isValid {
			issues = append(issues, problems...)
		}
	}

	var err error
	isValid := len(issues) < 1

	if !isValid {
		err = errors.New(strings.Join(issues, "\r\n"))
	}

	return isValid, err
}

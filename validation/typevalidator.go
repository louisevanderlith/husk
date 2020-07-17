package validation

import (
	"reflect"
)

//GetTypeValidator returns a Validation Function for a Specific Type
func GetTypeValidator(fieldType reflect.Kind) IValidation {
	var result IValidation

	switch fieldType {
	case reflect.Bool:
		result = typeValidator(happyValidation)
	case reflect.Int:
		result = typeValidator(intValidation)
	case reflect.Int8:
		result = typeValidator(int8Validation)
	case reflect.Int16:
		result = typeValidator(int16Validation)
	case reflect.Int32:
		result = typeValidator(intValidation)
	case reflect.Int64:
		result = typeValidator(int64Validation)
	case reflect.Uint:
		result = typeValidator(uintValidation)
	case reflect.Uint8:
		result = typeValidator(uint8Validation)
	case reflect.Uint16:
		result = typeValidator(uint16Validation)
	case reflect.Uint32:
		result = typeValidator(uintValidation)
	case reflect.Uint64:
		result = typeValidator(uint64Validation)
	case reflect.Float32:
		result = typeValidator(floatValidation)
	case reflect.Float64:
		result = typeValidator(float64Validation)
	case reflect.Complex64:
		result = typeValidator(happyValidation)
	case reflect.Complex128:
		result = typeValidator(happyValidation)
	case reflect.Array:
		result = typeValidator(happyValidation)
	case reflect.Map:
		result = typeValidator(happyValidation)
	case reflect.Ptr:
		result = typeValidator(pointerValidation)
	case reflect.Slice:
		result = typeValidator(happyValidation)
	case reflect.String:
		result = typeValidator(stringValidation)
	case reflect.Struct:
		result = typeValidator(structValidation)
	case reflect.Interface:
		result = typeValidator(interfaceValidation)
	default:
		result = typeValidator(noValidation)
	}

	return result
}

type typeValidator func(obj interface{}, meta tagMeta) (bool, []string)

func (t typeValidator) Valid(obj interface{}, meta tagMeta) (bool, []string) {
	return t(obj, meta)
}

package husk

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type tagMeta struct {
	Required bool
	Size     int
	Type     reflect.Kind
	PropName string
}

type IValidation interface {
	Valid(obj interface{}, meta tagMeta) (bool, []string)
}

const (
	idMessage       = "%s must be provided."
	emptyMessage    = "%s can't be empty."
	shortMessage    = "%s can't be more than %v characters."
	relationMessage = "%s can't be nil."
	incorrectType   = "%s's value '%s' is not of type %s."
)

// ValidateStruct will read 'hsk' tags on properties, to validate their values
// Properties without the 'hsk' tag, will be considered 'Required'
func ValidateStruct(obj interface{}) (bool, error) {
	var issues []string

	val := reflect.ValueOf(obj).Elem()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag.Get("hsk")

		kind := valueField.Kind()
		validator := getTypeValidator(kind)

		meta := getMeta(tag, kind)
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

func getIDMessage(property string) string {
	return fmt.Sprintf(idMessage, property)
}

func getEmptyMessage(property string) string {
	return fmt.Sprintf(emptyMessage, property)
}

func getShortMessage(property string, length int) string {
	return fmt.Sprintf(shortMessage, property, length)
}

func getRelationMessage(property string) string {
	return fmt.Sprintf(relationMessage, property)
}

func getIncorrectType(property string, value interface{}, correctType string) string {
	return fmt.Sprintf(incorrectType, property, value, correctType)
}

func getMeta(tag string, kind reflect.Kind) tagMeta {
	result := tagMeta{}
	parts := strings.Split(tag, ";")

	required := !strings.Contains(tag, "null")
	result.Required = required
	result.Type = kind

	hasSize := strings.Contains(tag, "size")

	if hasSize {
		rawSize := getFromTag(parts, "size")
		sSize := strings.Replace(strings.Replace(rawSize, "size(", "", -1), ")", "", -1)

		size, err := strconv.ParseInt(sSize, 10, 32)

		if err == nil {
			result.Size = int(size)
		}
	}

	return result
}

func getFromTag(list []string, name string) string {
	var result string
	for _, v := range list {
		if strings.Contains(v, name) {
			result = v
			break
		}
	}

	return result
}

func getTypeValidator(fieldType reflect.Kind) IValidation {
	var result IValidation

	switch fieldType {
	case reflect.Int:
		result = IntValidation{}
	case reflect.Int64:
		result = Int64Validation{}
	case reflect.String:
		result = StringValidation{}
	case reflect.Struct:
		result = StructValidation{}
	case reflect.Ptr:
		result = PointerValidation{}
	default:
		result = PointerValidation{}
	}

	return result
}

type StringValidation struct{}

func (o StringValidation) Valid(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(string)

	if ok {
		if meta.Required && val == "" {
			issues = append(issues, getEmptyMessage(meta.PropName))
		}

		if meta.Size > 0 && len(val) > meta.Size {
			issues = append(issues, getShortMessage(meta.PropName, meta.Size))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "string"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

type IntValidation struct{}

func (o IntValidation) Valid(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(int)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "int"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

type Int64Validation struct{}

func (o Int64Validation) Valid(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string
	val, ok := obj.(int64)

	if ok {
		if meta.Required && val < 1 {
			issues = append(issues, getIDMessage(meta.PropName))
		}
	} else {
		issues = append(issues, getIncorrectType(meta.PropName, obj, "int64"))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

type StructValidation struct{}

func (o StructValidation) Valid(obj interface{}, meta tagMeta) (bool, []string) {
	var issues []string

	if meta.Required && obj == nil {
		issues = append(issues, getRelationMessage(meta.PropName))
	}

	isValid := len(issues) < 1

	return isValid, issues
}

type PointerValidation struct{}

func (o PointerValidation) Valid(obj interface{}, meta tagMeta) (bool, []string) {
	return true, nil
}

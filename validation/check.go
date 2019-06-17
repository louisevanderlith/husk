package validation

import (
	"fmt"
)

//IValidation provides a Validation method for checking `hsk` tags
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

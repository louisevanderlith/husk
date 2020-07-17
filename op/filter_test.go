package op

import (
	"errors"
	"testing"
)

type aString string

func (s aString) Valid() (bool, error) {
	return false, errors.New("validation ran")
}

func TestDitto_NoError(t *testing.T) {
	aStringFilter := func(obj aString) bool {
		return obj == "PASS"
	}

	param := aString("NO PASS")
	pass := aStringFilter(param)

	if pass {
		t.Error("function returned true, unexepectedly")
	}
}

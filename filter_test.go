package husk

import (
	"errors"
	"testing"
)

type aString string

func (s aString) Valid() (bool, error) {
	return false, errors.New("validation ran")
}

type stringFilter func(obj aString) bool

func TestMakeFilter_NoError(t *testing.T) {
	var aStringFilter stringFilter
	aStringFilter = func(obj aString) bool {
		return obj == "PASS"
	}

	filter, err := MakeFilter(aStringFilter)

	if err != nil {
		t.Error(err)
	}

	param := aString("NO PASS")
	pass := filter(param)

	if pass {
		t.Error("function returned true, execpectedly")
	}
}

func TestMakeFilter_FuncOnly(t *testing.T) {

}

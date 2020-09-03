package validation

import (
	"html/template"
	"reflect"
	"testing"
)

type greatObject struct {
	Name string
}

type enum int

const (
	on enum = iota
	off
)

func (e enum) Valid() error {
	return nil
}

var vals = [...]string{
	"on",
	"off",
}

func TestGreatObject_Validation_OK(t *testing.T) {
	//in := greatObject{Name: ""}
}

func TestStringValidation_ValidString_OK(t *testing.T) {
	in := "A Real String"
	meta := GetMeta("", reflect.String)
	valid, err := stringValidation(in, meta)

	if len(err) > 0 {
		t.Error(err)
		return
	}

	if !valid {
		t.Error("Field has to be valid.")
	}
}

func TestStringValidation_HTMLString_OK(t *testing.T) {
	var in template.HTML
	in = "<p>A Real <b>String</b></p>"
	meta := GetMeta("", reflect.String)
	valid, err := stringValidation(in, meta)

	if len(err) > 0 {
		t.Error(err)
		return
	}

	if !valid {
		t.Error("Field has to be valid.")
	}
}

func TestEnumValidation_ValidInt_OK(t *testing.T) {
	meta := GetMeta("", reflect.Int)
	valid, err := interfaceValidation(on, meta)

	if len(err) > 0 {
		t.Error(err)
		return
	}

	if !valid {
		t.Error("Field has to be valid.")
	}
}

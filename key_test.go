package husk

import (
	"encoding/json"
	"testing"
)

func TestKey_CanParse(t *testing.T) {
	k := CrazyKey().String()
	prs, err := ParseKey(k)

	if err != nil {
		t.Error(err)
	}

	if prs.String() != k {
		t.Errorf("Expected %s, got %+v.", k, prs)
	}
}

func TestKey_TOJSON(t *testing.T) {
	k := CrazyKey()

	expected, _ := k.MarshalJSON()
	actual, err := json.Marshal(k)

	if err != nil {
		t.Error(err)
	}

	if string(actual) == string(expected) {
		t.Errorf("expected %s, got %s", string(expected), string(actual))
	}
}

func TestKey_FromJSON(t *testing.T) {
	expected := NewKey(2)
	input, _ := expected.MarshalJSON()

	actual := CrazyKey()
	err := json.Unmarshal(input, &actual)

	if err != nil {
		t.Error(err)
	}

	if actual.Stamp != expected.Stamp {
		t.Errorf("expected %v, got %v", expected.Stamp, actual.Stamp)
	}
}

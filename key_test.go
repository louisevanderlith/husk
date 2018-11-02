package husk

import (
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

package collections

import (
	"testing"
)

func TestComparable_NumCompare_15gt9(t *testing.T) {
	exp := int8(1) //larger
	in := compareNum(15)
	act := in.Compare(9)

	if act != exp {
		t.Error("expected", exp, "actual", act)
	}
}

func TestComparable_NumCompare_9lt15(t *testing.T) {
	exp := int8(-1) //smaller
	in := compareNum(9)
	act := in.Compare(15)

	if act != exp {
		t.Error("expected", exp, "actual", act)
	}
}

func TestCompareable_NumCompare_9eq9(t *testing.T) {
	exp := int8(0) //equal
	in := compareNum(9)
	act := in.Compare(9)

	if act != exp {
		t.Error("expected", exp, "actual", act)
	}
}

type compareNum int

//Compare returns -1 (smaller), 0 (equal), 1 (larger)
func (c compareNum) Compare(c2 compareNum) int8 {
	if c < c2 {
		return -1
	}

	if c > c2 {
		return 1
	}

	return 0
}

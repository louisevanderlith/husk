package records

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"testing"
)

type alpha string

func (a alpha) Valid() error {
	return nil
}

func TestList_NewList(t *testing.T) {
	l := NewCollection()

	if l.Count() != 0 {
		t.Error("count not zero")
	}
}

func TestList_Add_MustCount(t *testing.T) {
	l := NewCollection()
	rec := hsk.MakeRecord(keys.NewKey(1), alpha("A"))
	l.Add(rec)

	if l.Count() != 1 {
		t.Error("item count didn't increment", l.Count())
	}
}

func TestList_Insert_MustCount(t *testing.T) {
	l := NewCollection()
	rec := hsk.MakeRecord(keys.NewKey(1), alpha("A"))
	l.Insert(0, rec)

	if l.Count() != 1 {
		t.Error("item count didn't increment", l.Count())
	}
}

func TestList_Contains(t *testing.T) {
	l := NewCollection()
	rec := hsk.MakeRecord(keys.NewKey(1), alpha("A"))
	l.Add(rec)

	if !l.Contains(rec) {
		t.Error("unable to find item")
	}
}

func TestList_Get(t *testing.T) {
	l := NewCollection()
	xRec := hsk.MakeRecord(keys.NewKey(2), alpha("B"))
	l.Add(hsk.MakeRecord(keys.NewKey(1), alpha("A")))
	id := l.Add(xRec)
	l.Add(hsk.MakeRecord(keys.NewKey(3), alpha("C")))

	a := l.Get(id)

	if a != xRec {
		t.Error("unexpected value", a)
	}
}

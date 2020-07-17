package collections

import "testing"

func TestList_NewList(t *testing.T) {
	l := NewList()

	if l.Count() != 0 {
		t.Error("count not zero")
	}
}

func TestList_Add_MustCount(t *testing.T) {
	l := NewList()
	l.Add("A")

	if l.Count() != 1 {
		t.Error("item count didn't increment", l.Count())
	}
}

func TestList_Insert_MustCount(t *testing.T) {
	l := NewList()
	l.Insert(0, "A")

	if l.Count() != 1 {
		t.Error("item count didn't increment", l.Count())
	}
}

func TestList_Contains(t *testing.T) {
	l := NewList()
	l.Add("A")

	if !l.Contains("A") {
		t.Error("unable to find item")
	}
}

func TestList_Get(t *testing.T) {
	l := NewList()
	l.Add("A")
	l.Add("B")
	l.Add("C")

	a := l.Get(1)

	if a != "B" {
		t.Error("unexpected value", a)
	}
}

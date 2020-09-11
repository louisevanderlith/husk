package records

import (
	"encoding/json"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"testing"
)

type alpha struct {
	Char string
}

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
	rec := MakeRecord(keys.NewKey(1), &alpha{"B"})
	l.Add(rec)

	if l.Count() != 1 {
		t.Error("item count didn't increment", l.Count())
	}
}

func TestList_Insert_MustCount(t *testing.T) {
	l := NewCollection()
	rec := MakeRecord(keys.NewKey(1), alpha{"A"})
	l.Insert(0, rec)

	if l.Count() != 1 {
		t.Error("item count didn't increment", l.Count())
	}
}

func TestList_Contains(t *testing.T) {
	l := NewCollection()
	rec := MakeRecord(keys.NewKey(1), alpha{"B"})
	l.Add(rec)

	if !l.Contains(rec) {
		t.Error("unable to find item")
	}
}

func TestList_Get(t *testing.T) {
	l := NewCollection()
	xRec := MakeRecord(keys.NewKey(2), alpha{"B"})
	l.Add(MakeRecord(keys.NewKey(1), alpha{"A"}))
	id := l.Add(xRec)
	l.Add(MakeRecord(keys.NewKey(3), alpha{"C"}))

	a := l.Get(id)

	if a != xRec {
		t.Error("unexpected value", a)
	}
}

func TestList_MarshalJSON(t *testing.T) {
	in := NewCollection()
	in.Add(MakeRecord(keys.NewKeyWithTime(1599470402, 5), alpha{"XF"}))
	in.Add(MakeRecord(keys.NewKeyWithTime(1599470402, 6), alpha{"YT"}))
	act, err := json.Marshal(in)

	if err != nil {
		t.Error(err)
		return
	}

	exp := "[{\"Key\":\"1599470402`5\",\"Value\":{\"Char\":\"XF\"}},{\"Key\":\"1599470402`6\",\"Value\":{\"Char\":\"YT\"}}]"

	if string(act) != exp {
		t.Error("Expected", exp, "Got", string(act))
	}
}

func TestList_UnmarshalJSON(t *testing.T) {
	act := CollectionOf(alpha{})
	in := "[{\"Key\":\"1599470402`5\",\"Value\":{\"Char\":\"XF\"}},{\"Key\":\"1599470402`6\",\"Value\":{\"Char\":\"YT\"}}]"
	err := json.Unmarshal([]byte(in), act)

	if err != nil {
		t.Error(err)
		return
	}

	exp := NewCollection()
	exp.Add(MakeRecord(keys.NewKeyWithTime(1599470402, 5), alpha{"XF"}))
	exp.Add(MakeRecord(keys.NewKeyWithTime(1599470402, 6), alpha{"YT"}))

	if act.Count() != exp.Count() {
		t.Error("Count Expected", exp.Count(), "Got", act.Count())
	}

	itor := act.GetEnumerator()

	for itor.MoveNext() {
		rec := itor.Current().(hsk.Record)
		idx := exp.IndexOf(rec)

		if idx == -1 {
			t.Error(rec.GetKey(), rec.GetValue(), "not found")
		}
	}
}

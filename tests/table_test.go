package tests

import (
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/tests/sample"
	"testing"
)

func TestNewTable(t *testing.T) {
	tbl := husk.NewTable(sample.Event{})
	exp := "Event"
	if tbl.Name() != exp {
		t.Error("Invalid Name;", exp, "got", tbl.Name())
	}
}

func TestCreateAndFind(t *testing.T) {
	tbl := husk.NewTable(sample.Event{})

	in := sample.Event{Type: "CREATE", Relation: keys.CrazyKey()}
	k, err := tbl.Create(in)

	if err != nil {
		t.Fatal(err)
		return
	}

	rec, err := tbl.FindByKey(k)

	if err != nil {
		t.Error("Find Error", err)
		return
	}

	act := rec.GetValue().(sample.Event)
	if act.Type != in.Type {
		t.Error("Invalid Type; expected", in.Type, "got", act.Type)
	}
}

func TestUpdate(t *testing.T) {
	tbl := husk.NewTable(sample.Event{})

	in := sample.Event{Type: "CREATE", Relation: keys.CrazyKey()}
	k, err := tbl.Create(in)

	if err != nil {
		t.Fatal(err)
		return
	}

	in.Type = "UPDATE"
	err = tbl.Update(k, in)

	if err != nil {
		t.Fatal(err)
		return
	}

	rec, err := tbl.FindByKey(k)
	act := rec.GetValue().(sample.Event)
	if act.Type != in.Type {
		t.Error("Invalid Name; expected", in.Type, "got", in.Type)
	}
}

func TestDelete(t *testing.T) {
	tbl := husk.NewTable(sample.Event{})

	in := sample.Event{Type: "CREATE", Relation: keys.CrazyKey()}
	k, err := tbl.Create(in)

	if err != nil {
		t.Fatal(err)
		return
	}

	err = tbl.Delete(k)

	if err != nil {
		t.Fatal(err)
		return
	}

	_, err = tbl.FindByKey(k)

	if err == nil {
		t.Error("expected error")
	}
}

package tests

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/tests/sample"
	"testing"
)

/*
	Context tests should tests;
	1. Data is Saved to Disk
	2. Index is Saved to Disk
	3. Index & Data can be reloaded from Disk
	4. Seeding works, and saves to Disk.
*/

func TestUpdate_MustPersist(t *testing.T) {
	ctx := sample.NewEventContext()
	k, err := ctx.CreateEvent("INSERT", keys.CrazyKey())

	if err != nil {
		t.Error("Create Error", err)
		return
	}

	rec, err := ctx.GetEvent(k)

	if err != nil {
		t.Error(err)
		return
	}

	p := rec.GetValue().(sample.Event)
	p.Type = "DELETE"

	err = ctx.UpdateEvent(k, p)

	if err != nil {
		t.Error("Update Error", err)
		return
	}

	againP, err := ctx.GetEvent(k)

	if err != nil {
		t.Error("Get Error", err)
	}

	againData := againP.GetValue().(sample.Event)

	if againData.Type != p.Type {
		t.Errorf("Expected %v, got %v", p.Type, againData.Type)
	}
}

func TestDelete_MustPersist(t *testing.T) {
	ctx := sample.NewEventContext()
	k, err := ctx.CreateEvent("DELETE", keys.CrazyKey())

	if err != nil {
		t.Error(err)
		return
	}

	t.Log(k)

	err = ctx.DeleteEvent(k)

	if err != nil {
		t.Error(err)
		return
	}

	_, err = ctx.GetEvent(k)

	if err == nil {
		t.Error("Expected item to be deleted. 'Not found error...'")
	}
}

func TestFind_FindFilteredItems(t *testing.T) {
	ctx := sample.NewEventContext()

	ka, err := ctx.CreateEvent("INSERT", keys.CrazyKey())

	if err != nil {
		t.Fatal(err)
	}
	kb, err := ctx.CreateEvent("READ", keys.CrazyKey())

	if err != nil {
		t.Fatal(err)
	}
	kc, err := ctx.CreateEvent("DELETE", keys.CrazyKey())

	if err != nil {
		t.Fatal(err)
	}

	result, err := ctx.FindEventsByType(1, 5, "INSERT")

	if err != nil {
		t.Error(err)
		return
	}

	itor := result.GetEnumerator()
	matchFound := false
	for itor.MoveNext() {
		curr := itor.Current().(hsk.Record)
		recKey := curr.GetKey()

		for _, k := range []hsk.Key{ka, kb, kc} {
			if recKey == k {
				matchFound = true
				break
			}
		}
	}

	if !matchFound {
		t.Error("no matches found")
	}
}

func TestFind_FindEverything(t *testing.T) {
	ctx := sample.NewEventContext()
	rset, err := ctx.FindEvents(1, 100)

	if err != nil {
		t.Fatal(err)
		return
	}

	if !rset.Any() {
		t.Error("no data found")
	}
}

func TestFilter_FindEverything_MustBe10(t *testing.T) {
	ctx := sample.NewEventContext()
	records, err := ctx.FindEvents(1, 10)

	if err != nil {
		t.Fatal(err)
		return
	}

	if records.Count() != 10 {
		t.Errorf("Expecting 10, got %d", records.Count())
	}
}

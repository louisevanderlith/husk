package tests

import (
	"encoding/json"
	"github.com/louisevanderlith/husk/db"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"strings"
	"testing"
)

/*
	Context tests should tests;
	1. Data is Saved to Disk
	2. Index is Saved to Disk
	3. Index & Data can be reloaded from Disk
	4. Seeding works, and saves to Disk.
*/

func TestDiscover_ListNames(t *testing.T) {
	ctx := db.NewContext()
	exp := []string{"Events"}
	act := db.TableNames(ctx)

	if len(act) != len(exp) {
		t.Error("invalid length discovered")
		return
	}

	if exp[0] != act[0] {
		t.Errorf("Expected %s, found %s", exp, act)
	}
}

func TestDiscover_ListLayouts(t *testing.T) {
	ctx := db.NewContext()
	act := db.TableLayouts(ctx)

	if len(act) != 1 {
		t.Error("invalid length discovered")
		return
	}

	if act["Events"] == nil {
		t.Errorf("no object found %v", act)
	}
}

func TestCount_JournalCount(t *testing.T) {
	count, err := benchCtx.CountJournals()

	if err != nil {
		t.Fatal(err)
		return
	}

	if count == 0 {
		t.Fatal("invalid count")
	}
}

func TestFind_SearchItems(t *testing.T) {
	set, err := benchCtx.FindJournalsByPublisher(1, 10, "University of Malaya")

	if err != nil {
		t.Error(err)
		return
	}

	itor := set.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current().(hsk.Record)

		t.Log(curr.Data())
	}

	if set.Count() != 5 {
		t.Errorf("%+v\n", set.Count())
	}
}

func TestUpdate_MustPersist(t *testing.T) {
	p := db.Event{Type: "INSERT", RecordKey: keys.CrazyKey()}

	ctx := db.NewContext()
	k, err := ctx.CreateEvent(p)

	if err != nil {
		t.Error("Create Error", err)
		return
	}

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

	againData := againP.Data().(db.Event)

	if againData.Type != p.Type {
		t.Errorf("Expected %v, got %v", p.Type, againData.Type)
	}
}

func TestDelete_MustPersist(t *testing.T) {
	p := db.Event{Type: "DELETE", RecordKey: keys.CrazyKey()}

	ctx := db.NewContext()
	k, err := ctx.CreateEvent(p)

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
	p := db.Event{Type: "INSERT", RecordKey: keys.CrazyKey()}
	p1 := db.Event{Type: "READ", RecordKey: keys.CrazyKey()}
	p2 := db.Event{Type: "DELETE", RecordKey: keys.CrazyKey()}

	ctx := db.NewContext()
	ks, err := ctx.CreateEvents(p, p1, p2)

	if err != nil {
		t.Fatal(err)
	}

	result, err := ctx.ListEventsWithType("INSERT")

	if err != nil {
		t.Error(err)
		return
	}

	itor := result.GetEnumerator()
	matchFound := false
	for itor.MoveNext() {
		curr := itor.Current().(hsk.Record)
		recKey := curr.GetKey()

		for _, k := range ks {
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
	ctx := db.NewContext()
	rset, err := ctx.ListEvents()

	if err != nil {
		t.Fatal(err)
		return
	}

	if !rset.Any() {
		t.Error("no data found")
	}
}

func TestFilter_FindEverything_MustBe10(t *testing.T) {
	ctx := db.NewContext()
	records, err := ctx.ListEvents()

	if err != nil {
		t.Fatal(err)
		return
	}

	if records.Count() != 10 {
		t.Errorf("Expecting 10, got %d", records.Count())
	}
}

func TestRecordSet_ToJSON_MustBeClean(t *testing.T) {
	ctx := db.NewContext()
	rows, err := ctx.ListEvents()

	if err != nil {
		t.Fatal(err)
		return
	}

	bits, _ := json.Marshal(rows)

	if strings.Contains(string(bits), "Value") {
		t.Error("Final Object has Value")
	}
}

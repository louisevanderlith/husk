package tests

import (
	"encoding/json"
	"github.com/louisevanderlith/husk/db"
	"github.com/louisevanderlith/husk/hsk"
	"strings"
	"testing"
)

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
		curr := itor.Current().(hsk.Recorder)

		t.Log(curr.Data())
	}

	if set.Count() != 6 {
		t.Errorf("%+v\n", set.Count())
	}
}

func TestUpdate_MustPersist(t *testing.T) {
	p := db.Event{Type: "INSERT", RecordKey: hsk.CrazyKey()}

	ctx := db.NewContext()
	k, err := ctx.CreateEvent(p)

	if err != nil {
		t.Error(err)
		return
	}

	p.Type = "DELETE"

	err = ctx.UpdateEvent(k, p)

	if err != nil {
		t.Error(err)
		return
	}

	againP, err := ctx.GetEvent(k)

	if err != nil {
		t.Error(err)
	}

	againData := againP.Data().(db.Event)

	if againData.Type != p.Type {
		t.Errorf("Expected %v, got %v", p.Type, againData.Type)
	}
}

func TestDelete_MustPersist(t *testing.T) {
	p := db.Event{Type: "DELETE", RecordKey: hsk.CrazyKey()}

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
	p := db.Event{Type: "INSERT", RecordKey: hsk.CrazyKey()}
	p1 := db.Event{Type: "READ", RecordKey: hsk.CrazyKey()}
	p2 := db.Event{Type: "DELETE", RecordKey: hsk.CrazyKey()}

	ctx := db.NewContext()
	keys, err := ctx.CreateEvents(p, p1, p2)

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
		curr := itor.Current().(hsk.Recorder)
		recKey := curr.GetKey()

		for _, k := range keys {
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

func TestFilter_FindEverything_MustBe1000(t *testing.T) {
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

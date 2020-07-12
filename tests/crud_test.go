package tests

import (
	"github.com/louisevanderlith/husk/db"
	"github.com/louisevanderlith/husk/hsk"
	"testing"
)

func TestCreate_MustPersist(t *testing.T) {
	ctx := db.NewContext()

	p := db.Event{
		Type:      "INSERT",
		RecordKey: hsk.CrazyKey(),
	}

	k, err := ctx.CreateEvent(p)

	if err != nil {
		t.Error(err)
	}

	againP, err := ctx.GetEvent(k)

	if err != nil {
		t.Error("Find Error", err)
		return
	}

	if againP == nil {
		t.Error("Record not found")
		return
	}

	if againP.GetKey() != k {
		t.Errorf("Expected %s, %s", k, againP.GetKey())
	}
}

func TestCreate_MultipleEntries_MustPersist(t *testing.T) {
	p := db.Event{Type: "INSERT", RecordKey: hsk.CrazyKey()}
	p1 := db.Event{Type: "READ", RecordKey: hsk.CrazyKey()}
	p2 := db.Event{Type: "DELETE", RecordKey: hsk.CrazyKey()}

	ctx := db.NewContext()
	keys, err := ctx.CreateEvents(p, p1, p2)

	if err != nil {
		t.Error(err)
		return
	}

	if len(keys) != 3 {
		t.Error("expected", 3, "got", len(keys))
		return
	}

	for _, k := range keys {
		_, err := ctx.GetEvent(k)

		if err != nil {
			t.Error(err)
			return
		}
	}
}



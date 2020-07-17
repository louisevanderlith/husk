package tests

import (
	"github.com/louisevanderlith/husk/db"
	"github.com/louisevanderlith/husk/keys"
	"testing"
)

func TestCreate_MustPersist(t *testing.T) {
	ctx := db.NewContext()

	p := db.Event{
		Type:      "INSERT",
		RecordKey: keys.CrazyKey(),
	}

	k, err := ctx.CreateEvent(p)

	if err != nil {
		t.Error(err)
		return
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
	p := db.Event{Type: "INSERT", RecordKey: keys.CrazyKey()}
	p1 := db.Event{Type: "READ", RecordKey: keys.CrazyKey()}
	p2 := db.Event{Type: "DELETE", RecordKey: keys.CrazyKey()}

	ctx := db.NewContext()
	ks, err := ctx.CreateEvents(p, p1, p2)

	if err != nil {
		t.Error(err)
		return
	}

	if len(ks) != 3 {
		t.Error("expected", 3, "got", len(ks))
		return
	}

	for _, k := range ks {
		_, err := ctx.GetEvent(k)

		if err != nil {
			t.Error(err)
			return
		}
	}
}

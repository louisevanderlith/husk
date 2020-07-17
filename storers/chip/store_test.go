package chip

import (
	"github.com/louisevanderlith/husk/db"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"testing"
)

func TestChipStore_Write(t *testing.T) {
	store := newStore()

	in := db.Event{
		Type:      "WRITE",
		RecordKey: keys.CrazyKey(),
	}

	p, err := store.Write(in)

	if err != nil {
		t.Error("Write Error", err)
		return
	}

	if p.GetLength() == 0 {
		t.Error("Unexpected Length", p.GetLength())
		return
	}
}

func TestChipStore_Read(t *testing.T) {
	store := newStore()

	in := db.Event{
		Type:      "WRITE",
		RecordKey: keys.CrazyKey(),
	}

	p, err := store.Write(in)

	if err != nil {
		t.Error("Write Error", err)
		return
	}

	data := make(chan hsk.Dataer)
	err = store.Read(p, data)

	if err != nil {
		t.Error("Read Error", err)
		return
	}

	obj := (<-data).(db.Event)
	if obj.Type != in.Type {
		t.Error("Invalid Type; expected", in.Type, "got", obj.Type)
		return
	}

	if obj.RecordKey != in.RecordKey {
		t.Error("Invalid RecordKey; expected", in.RecordKey, "got", obj.RecordKey)
		return
	}
}

func TestChipStore_Close(t *testing.T) {
	store := newStore()
	err := store.Close()

	if err != nil {
		t.Error("Close Error", err)
		return
	}
}

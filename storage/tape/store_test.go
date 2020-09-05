package tape

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/persisted"
	"github.com/louisevanderlith/husk/validation"
	"testing"
)

type tapeRecord struct {
	Name string `hsk:"size(128)"`
}

func (r tapeRecord) Valid() error {
	return validation.Struct(r)
}

func BenchmarkTapeStore_Write(b *testing.B) {
	store := NewDefaultStore(tapeRecord{})

	in := tapeRecord{"WRITE"}
	rec := hsk.MakeRecord(keys.CrazyKey(), in)
	p := make(chan hsk.Point)
	go store.Write(rec, p)

	res := <-p
	if res.GetLength() == 0 {
		b.Error("Unexpected Length", res.GetLength())
		return
	}
}

func TestTapeStore_Write(t *testing.T) {
	persisted.CreateDirectory("db")
	store := NewDefaultStore(tapeRecord{})

	in := tapeRecord{"WRITE"}
	rec := hsk.MakeRecord(keys.CrazyKey(), in)
	p := make(chan hsk.Point)
	go store.Write(rec, p)

	res := <-p
	if res.GetLength() == 0 {
		t.Error("Unexpected Length", res.GetLength())
		return
	}
}

func BenchmarkTapeStore_Read(b *testing.B) {
	store := NewDefaultStore(tapeRecord{})

	in := tapeRecord{"WRITE"}
	rec := hsk.MakeRecord(keys.CrazyKey(), in)
	p := make(chan hsk.Point)
	go store.Write(rec, p)

	data := make(chan hsk.Record)
	go store.Read(<-p, data)

	obj := (<-data).Data().(tapeRecord)

	if obj.Name != in.Name {
		b.Error("Invalid Name; expected", in.Name, "got", obj.Name)
		return
	}
}

func TestTapeStore_Read(t *testing.T) {
	store := NewDefaultStore(tapeRecord{})

	in := tapeRecord{"WRITE"}
	rec := hsk.MakeRecord(keys.CrazyKey(), in)
	p := make(chan hsk.Point)
	go store.Write(rec, p)

	data := make(chan hsk.Record)
	go store.Read(<-p, data)

	obj := (<-data).Data().(tapeRecord)
	if obj.Name != in.Name {
		t.Error("Invalid Name; expected", in.Name, "got", obj.Name)
		return
	}
}

func TestTapeStore_Close(t *testing.T) {
	store := NewDefaultStore(tapeRecord{})
	err := store.Close()

	if err != nil {
		t.Error("Close Error", err)
		return
	}
}

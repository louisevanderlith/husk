package tape

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/validation"
	"reflect"
	"testing"
)

type tapeRecord struct {
	Name string `hsk:"size(128)"`
}

func (r tapeRecord) Valid() error {
	return validation.ValidateStruct(r)
}

func TestTapeStore_Write(t *testing.T) {
	store := NewStore(reflect.TypeOf(tapeRecord{}))

	in := tapeRecord{"WRITE"}

	p := make(chan hsk.Point)
	go store.Write(in, p)

	res := <-p
	if res.GetLength() == 0 {
		t.Error("Unexpected Length", res.GetLength())
		return
	}
}

func TestTapeStore_Read(t *testing.T) {
	store := NewStore(reflect.TypeOf(tapeRecord{}))

	in := tapeRecord{"WRITE"}

	p := make(chan hsk.Point)
	go store.Write(in, p)

	data := make(chan validation.Dataer)
	go store.Read(<-p, data)

	obj := (<-data).(tapeRecord)
	if obj.Name != in.Name {
		t.Error("Invalid Name; expected", in.Name, "got", obj.Name)
		return
	}
}

func TestTapeStore_Close(t *testing.T) {
	store := NewStore(reflect.TypeOf(tapeRecord{}))
	err := store.Close()

	if err != nil {
		t.Error("Close Error", err)
		return
	}
}

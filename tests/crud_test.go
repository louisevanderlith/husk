package tests

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/tests/sample"
	"testing"
)

func TestCreate_MustPersist(t *testing.T) {
	ctx := sample.NewEventContext()

	k, err := ctx.CreateEvent("INSERT", keys.CrazyKey())

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
	ctx := sample.NewEventContext()
	ka, err := ctx.CreateEvent("INSERT", keys.CrazyKey())

	if err != nil {
		t.Fatal(err)
		return
	}
	kb, err := ctx.CreateEvent("READ", keys.CrazyKey())

	if err != nil {
		t.Fatal(err)
		return
	}
	kc, err := ctx.CreateEvent("DELETE", keys.CrazyKey())

	if err != nil {
		t.Fatal(err)
		return
	}

	for _, k := range []hsk.Key{ka, kb, kc} {
		_, err := ctx.GetEvent(k)

		if err != nil {
			t.Error(err)
			return
		}
	}
}

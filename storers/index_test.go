package storers

import (
	"github.com/louisevanderlith/husk/hsk"
	"testing"
)

func TestNewIndex_Add_GetsKey(t *testing.T) {
	idx := NewIndex()
	mta := idx.Add(hsk.NewPoint(0, 1))

	if mta.GetKey() == hsk.CrazyKey() {
		t.Error("crazy key not expected")
	}
}

func TestNewIndex_Add_GetsMeta(t *testing.T) {
	idx := NewIndex()
	mta := idx.Add(hsk.NewPoint(0, 1))

	if mta.GetKey() == hsk.CrazyKey() {
		t.Error("crazy key not expected")
	}

	itm := idx.Get(mta.GetKey())

	if itm.Point().Offset != 0 || itm.Point().Len != 1 {
		t.Error("incorrect value found")
	}
}

func TestNewIndex_AddMany_GetsAllKeys(t *testing.T) {
	idx := NewIndex()

	for i := 0; i < 10; i++ {
		mta := idx.Add(hsk.NewPoint(0, 1))

		if mta.GetKey() == hsk.CrazyKey() {
			t.Error("crazy key not expected")
		}
	}

	if len(idx.GetKeys()) != 10 {
		t.Error("didn't add all items")
	}
}

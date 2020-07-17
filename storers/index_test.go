package storers

import (
	"github.com/louisevanderlith/husk/hsk"
	"testing"
)

func TestNewIndex_Add_GetsKey(t *testing.T) {
	idx := NewIndex()
	mta := hsk.NewMeta(hsk.NewPoint(0, 1))
	k, err := idx.Add(mta)

	if err != nil {
		t.Error("Add Error", err)
		return
	}

	if k == nil {
		t.Error("key is nil")
	}
}

func TestNewIndex_Add_GetsMeta(t *testing.T) {
	idx := NewIndex()
	mta := hsk.NewMeta(hsk.NewPoint(0, 1))
	k, err := idx.Add(mta)

	if err != nil {
		t.Error("Add Error", err)
		return
	}

	if k == nil {
		t.Error("key is nil")
		return
	}

	itm := idx.Get(k)

	itmPoint := itm.Point()
	mtaPoint := mta.Point()

	if itmPoint.GetOffset() != mtaPoint.GetOffset() {
		t.Error("Point Offset invalid. Expected", itmPoint.GetLength(), "got", mtaPoint.GetLength())
		return
	}

	if itmPoint.GetLength() != mtaPoint.GetLength() {
		t.Error("Point Length invalid. Expected", itmPoint.GetLength(), "got", mtaPoint.GetLength())
		return
	}
}

func TestNewIndex_AddMany_GetsAllKeys(t *testing.T) {
	idx := NewIndex()

	for i := int64(0); i < 10; i++ {
		mta := hsk.NewMeta(hsk.NewPoint(i, i+1))
		k, err := idx.Add(mta)

		if err != nil {
			t.Error("Add Error", err)
			return
		}

		if k == nil {
			t.Error("key is nil")
			return
		}
	}

	if len(idx.GetKeys()) != 10 {
		t.Error("didn't add all items")
	}
}

func TestIndex_IndexOf(t *testing.T) {
	idx := NewIndex()
	mta := hsk.NewMeta(hsk.NewPoint(0, 1))
	k, err := idx.Add(mta)

	if err != nil {
		t.Error("Add Error", err)
		return
	}

	gmta := idx.Get(k)

	if gmta == nil {
		t.Error("Meta is nil")
	}

	if gmta.Point().GetLength() != mta.Point().GetLength() {
		t.Error("expected", mta.Point().GetLength(), "got", gmta.Point().GetLength())
	}
}

func TestSaveIndex(t *testing.T) {
	DestroyContents("db")
	idx := NewIndex()

	for i := int64(0); i < 10; i++ {
		mta := hsk.NewMeta(hsk.NewPoint(i, i+1))
		k, err := idx.Add(mta)

		if err != nil {
			t.Error("Add Error", err)
			return
		}

		if k == nil {
			t.Error("key is nil")
			return
		}
	}

	err := saveIndex("Event", idx)

	if err != nil {
		t.Error("Save Index Error", err)
		return
	}
}

func TestSaveIndex_Reload(t *testing.T) {
	DestroyContents("db")
	idx := NewIndex()

	for i := int64(0); i < 10; i++ {
		mta := hsk.NewMeta(hsk.NewPoint(i, i+1))
		k, err := idx.Add(mta)

		if err != nil {
			t.Error("Add Error", err)
			return
		}

		if k == nil {
			t.Error("key is nil")
			return
		}
	}

	err := saveIndex("Event", idx)

	if err != nil {
		t.Error("Save Index Error", err)
		return
	}

	idxFile, err := openIndex("Event")

	if err != nil {
		t.Error("Open Index Error", err)
		return
	}

	sidx, err := loadIndex(idxFile)

	if err != nil {
		t.Error("Load Index Error", err)
		return
	}

	for _, k := range sidx.GetKeys() {
		meta := idx.Get(k)

		if meta == nil {
			t.Error("unable to find", k)
			return
		}

		if !meta.IsActive() {
			t.Error("record should be active")
			return
		}
	}
}

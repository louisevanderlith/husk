package tests

import (
	"github.com/louisevanderlith/husk/db"
	"github.com/louisevanderlith/husk/hsk"
	"testing"
)

func TestNext_ShouldReturnNext(t *testing.T) {
	ctx := db.NewContext()
	results, err := ctx.ListEvents()

	if err != nil {
		t.Fatal(err)
		return
	}

	if results.Count() == 0 {
		t.Error("no results")
		return
	}

	rator := results.GetEnumerator()
	for rator.MoveNext() {
		curr := rator.Current().(hsk.Recorder)
		t.Logf("%+v", curr.Data())
	}
}

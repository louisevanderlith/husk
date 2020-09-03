package tests

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/tests/sample"
	"testing"
)

func TestNext_ShouldReturnNext(t *testing.T) {
	ctx := sample.NewEventContext()
	results, err := ctx.FindEvents(1, 100)

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
		curr := rator.Current().(hsk.Record)
		t.Logf("%+v", curr.Data())
	}
}

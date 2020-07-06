package tests

import (
	"testing"

	"github.com/louisevanderlith/husk"
)

func TestNext_ShouldReturnNext(t *testing.T) {
	defer DestroyData()
	results, err := ctx.Journals.Find(1, 3, husk.Everything())

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
		curr := rator.Current()
		t.Logf("%+v", curr.Data())
	}
}

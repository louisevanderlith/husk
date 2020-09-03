package tests

import (
	"github.com/louisevanderlith/husk/hsk"
	"testing"

	"github.com/louisevanderlith/husk/tests/sample"
)

var (
	benchCtx sample.JournalContext
)

func BenchmarkExist_Everything(b *testing.B) {
	benchCtx = sample.NewContext()
	benchCtx.HasJournals()
}

func BenchmarkCount_JournalCount(b *testing.B) {
	count, err := benchCtx.CountJournals()

	if err != nil {
		b.Fatal(err)
		return
	}

	b.Log(count)
}

func BenchmarkFilter_FindByAuthor(b *testing.B) {
	page, err := benchCtx.FindJournalsByPublisher(1, 10, "Universidade Federal do Rio Grande")

	if err != nil {
		b.Error(err)
		return
	}

	itor := page.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current().(hsk.Record)

		b.Log(curr.Data())
	}

	if page.Count() != 6 {
		b.Errorf("%+v\n", page.Count())
	}
}

func TestCount_JournalCount(t *testing.T) {
	count, err := benchCtx.CountJournals()

	if err != nil {
		t.Fatal(err)
		return
	}

	if count == 0 {
		t.Fatal("invalid count")
	}
}

func TestFind_SearchItems(t *testing.T) {
	set, err := benchCtx.FindJournalsByPublisher(1, 10, "University of Malaya")

	if err != nil {
		t.Error(err)
		return
	}

	itor := set.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current().(hsk.Record)

		t.Log(curr.Data())
	}

	if set.Count() != 5 {
		t.Errorf("%+v\n", set.Count())
	}
}

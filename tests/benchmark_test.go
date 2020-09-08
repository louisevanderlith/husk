package tests

import (
	"github.com/louisevanderlith/husk/hsk"
	"testing"

	"github.com/louisevanderlith/husk/tests/sample"
)

func BenchmarkExist_Everything(b *testing.B) {
	benchCtx := sample.NewContext()
	benchCtx.HasJournals()
}

func BenchmarkCount_JournalCount(b *testing.B) {
	benchCtx := sample.NewContext()
	count, err := benchCtx.CountJournals()

	if err != nil {
		b.Fatal(err)
		return
	}

	b.Log(count)
}

func BenchmarkFilter_FindByAuthor(b *testing.B) {
	benchCtx := sample.NewContext()
	page, err := benchCtx.FindJournalsByPublisher(1, 10, "Universidade Federal do Rio Grande")

	if err != nil {
		b.Error(err)
		return
	}

	itor := page.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current().(hsk.Record)

		b.Log(curr.GetValue())
	}

	if page.Count() != 6 {
		b.Errorf("%+v\n", page.Count())
	}
}

func TestCount_JournalCount(t *testing.T) {
	benchCtx := sample.NewContext()
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
	benchCtx := sample.NewContext()
	set, err := benchCtx.FindJournalsByPublisher(1, 10, "University of Malaya")

	if err != nil {
		t.Error(err)
		return
	}

	itor := set.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current().(hsk.Record)

		t.Log(curr.GetValue())
	}

	if set.Count() != 5 {
		t.Errorf("%+v\n", set.Count())
	}
}

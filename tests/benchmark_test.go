package tests

import (
	"github.com/louisevanderlith/husk/hsk"
	"testing"

	"github.com/louisevanderlith/husk/tests/sample"
)

var (
	benchCtx sample.SampleContext
)

func init() {
	benchCtx = sample.NewContext()
}

func BenchmarkExist_Everything(b *testing.B) {
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

package tests

import (
	"github.com/louisevanderlith/husk"
	"testing"

	"github.com/louisevanderlith/husk/tests/sample"
)

var (
	benchCtx sample.Context
)

func init() {
	benchCtx = sample.NewContext()
	benchCtx.Seed()
}

func BenchmarkExist_Everything(b *testing.B) {
	benchCtx.Journals.Exists(husk.Everything())
}

func BenchmarkCount_JournalCount(b *testing.B) {
	count := int64(0)
	err := benchCtx.Journals.Calculate(&count, husk.RowCount())

	if err != nil {
		b.Fatal(err)
		return
	}

	b.Log(count)
}

func BenchmarkFilter_FindByAuthor(b *testing.B) {
	//defer DestroyData()
	set, err := benchCtx.Journals.Find(1, 10, sample.ByPublisher("Universidade Federal do Rio Grande"))

	if err != nil {
		b.Error(err)
		return
	}

	itor := set.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current()

		b.Log(curr.Data())
	}

	if set.Count() != 6 {
		b.Errorf("%+v\n", set.Count())
	}
}

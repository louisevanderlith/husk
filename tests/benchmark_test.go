package tests

import (
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

func BenchmarkFilter_PerfectBalance(b *testing.B) {
	set := benchCtx.People.Find(1, 10, sample.SameBalance(30321.12))

	itor := set.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current()

		b.Log(curr.Data())
	}

	if set.Count() != 10 {
		b.Errorf("%+v\n", set.Count())
	}
}

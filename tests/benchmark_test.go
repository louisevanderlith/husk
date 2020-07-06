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

}

func BenchmarkSeed_PeoplePopulate(b *testing.B) {
	DestroyData()
	err := benchCtx.People.Seed("people.seed.json")

	if err != nil {
		b.Fatal(err)
		return
	}

	benchCtx.People.Save()
}

func BenchmarkFilter_PerfectBalance(b *testing.B) {
	set, err := benchCtx.People.Find(1, 10, sample.SameBalance(30321.12))

	if err != nil {
		b.Error(err)
		return
	}

	itor := set.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current()

		b.Log(curr.Data())
	}

	if set.Count() != 10 {
		b.Errorf("%+v\n", set.Count())
	}
}

func BenchmarkCalc_CalculateLowestBalance(b *testing.B) {
	name := ""

	ctx.People.Calculate(&name, sample.LowestBalance())

	if name != "Kelley" {
		b.Errorf("name Kelley not found, got %s", name)
	}
}

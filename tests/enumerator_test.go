package tests

import (
	"testing"

	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/tests/sample"
)

func TestNext_ShouldReturnNext(t *testing.T) {
	defer DestroyData()
	names := []string{"Johan", "Sarel", "Jaco"}
	p := sample.Person{Name: names[0], Age: 13}
	ctx.People.Create(p)
	defer ctx.People.Save()

	p1 := sample.Person{Name: names[1], Age: 15}
	ctx.People.Create(p1)

	p2 := sample.Person{Name: names[2], Age: 24}
	ctx.People.Create(p2)

	results := ctx.People.Find(1, 3, husk.Everything())

	if results.Count() == 0 {
		t.Error("no results")
	}

	rator := results.GetEnumerator()
	for rator.MoveNext() {
		curr := rator.Current()
		t.Logf("%+v", curr.Data())
	}
}

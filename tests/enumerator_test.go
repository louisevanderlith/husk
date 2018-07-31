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

	p1 := sample.Person{Name: names[1], Age: 15}
	ctx.People.Create(p1)

	p2 := sample.Person{Name: names[2], Age: 24}
	ctx.People.Create(p2)

	results := ctx.People.Find(1, 3, func(obj husk.Dataer) bool {
		return true
	})

	for i := 0; i < 3; i++ {
		name := names[i]
		item, err := results.Next()

		if err != nil {
			t.Error(err)
		}

		per0 := item.Data().(*sample.Person)

		t.Log(per0.Name, name)
	}

	t.Fail()
}

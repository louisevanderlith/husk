package tests

import (
	"log"
	"testing"

	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/tests/sample"
)

var ctx sample.Context

func init() {
	ctx = sample.NewContext()
}

func DestroyData() {
	err := husk.DestroyContents("db")

	if err != nil {
		log.Println(err)
	}
}

func TestCreate_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Jan", Age: 25}

	set := ctx.People.Create(&p)
	defer ctx.People.Save()

	if set.Error != nil {
		t.Error(set.Error)
	}

	againP, ferr := ctx.People.FindByKey(set.Record.GetKey())

	if ferr != nil {
		t.Error(ferr)
	}

	if againP == nil {
		t.Error("Record not found")
	}
}

func TestCreate_MultipleEntries_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Johan", Age: 13}
	p1 := sample.Person{Name: "Sarel", Age: 15}
	p2 := sample.Person{Name: "Jaco", Age: 24}

	ctx.People.Create(p)
	ctx.People.Create(p1)
	p2Set := ctx.People.Create(p2)

	if p2Set.Error != nil {
		t.Error(p2Set.Error)
	}

	_, err := ctx.People.FindByKey(p2Set.Record.GetKey())

	if err != nil {
		t.Error(err)
	}
}

func TestUpdate_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Sarie", Age: 45}

	set := ctx.People.Create(&p)

	if set.Error != nil {
		t.Error(set)
	}

	pData := set.Record.Data().(*sample.Person)
	pData.Age = 67

	err := ctx.People.Update(set.Record)

	if err != nil {
		t.Error(err)
	}

	againP, ferr := ctx.People.FindByKey(set.Record.GetKey())

	if ferr != nil {
		t.Error(ferr)
	}

	againData := againP.Data().(*sample.Person)

	if againData.Age != p.Age {
		t.Errorf("Expected %v, got %v", p.Age, againData.Age)
	}
}

func TestUpdate_LastUpdatedMustChange(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Sarie", Age: 45}

	set := ctx.People.Create(&p)

	if set.Error != nil {
		t.Error(set.Record)
	}

	firstUpdate := set.Record.Meta().LastUpdated

	pData := set.Record.Data().(*sample.Person)
	pData.Age = 67

	err := ctx.People.Update(set.Record)

	if err != nil {
		t.Error(err)
	}

	againP, ferr := ctx.People.FindByKey(set.Record.GetKey())

	if ferr != nil {
		t.Error(ferr)
	}

	againMeta := againP.Meta()

	if againMeta.LastUpdated == firstUpdate {
		t.Errorf("Expected %v, got %v", firstUpdate, againMeta.LastUpdated)
	}
}

func TestDelete_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "DeleteMe", Age: 67}

	set := ctx.People.Create(p)

	if set.Error != nil {
		t.Error(set)
	}

	err := ctx.People.Delete(set.Record.GetKey())

	if err != nil {
		t.Error(err)
	}

	_, rerr := ctx.People.FindByKey(set.Record.GetKey())

	if rerr == nil {
		t.Error("Expected item to be deleted. 'Not found error...'")
	}
}

func TestFind_FindFilteredItems(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Johan", Age: 13}
	p1 := sample.Person{Name: "Sarel", Age: 15}
	p2 := sample.Person{Name: "Jaco", Age: 24}

	ctx.People.Create(p)
	set := ctx.People.Create(p1)
	ctx.People.Create(p2)

	result := ctx.People.FindFirst(func(obj husk.Dataer) bool {
		return obj.(*sample.Person).Name == "Sarel"
	})

	firstID := result.GetKey()

	if firstID != set.Record.GetKey() {
		t.Errorf("Wrong ID, Expected %v, got %v", set.Record.GetKey(), firstID)
	}
}

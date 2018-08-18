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
	//defer DestroyData()

	p := sample.Person{Name: "Jan", Age: 25}

	record, err := ctx.People.Create(&p)

	if err != nil {
		t.Error(err)
	}

	againP, ferr := ctx.People.FindByKey(record.GetKey())

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
	p2Record, rerr := ctx.People.Create(p2)

	if rerr != nil {
		t.Error(rerr)
	}

	_, err := ctx.People.FindByKey(p2Record.GetKey())

	if err != nil {
		t.Error(err)
	}
}

func TestUpdate_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Sarie", Age: 45}

	record, err := ctx.People.Create(&p)

	if err != nil {
		t.Error(err)
	}

	pData := record.Data().(*sample.Person)
	pData.Age = 67

	err = ctx.People.Update(record)

	if err != nil {
		t.Error(err)
	}

	againP, ferr := ctx.People.FindByKey(record.GetKey())

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

	record, err := ctx.People.Create(&p)

	if err != nil {
		t.Error(err)
	}

	firstUpdate := record.Meta().LastUpdated

	pData := record.Data().(*sample.Person)
	pData.Age = 67

	err = ctx.People.Update(record)

	if err != nil {
		t.Error(err)
	}

	againP, ferr := ctx.People.FindByKey(record.GetKey())

	if ferr != nil {
		t.Error(ferr)
	}

	againMeta := againP.Meta()

	if againMeta.LastUpdated == firstUpdate {
		t.Errorf("Expected %v, got %v", firstUpdate, againMeta.LastUpdated)
	}
}

func TestDelete_MustPersist(t *testing.T) {
	//defer DestroyData()

	p := sample.Person{Name: "DeleteMe", Age: 67}

	record, err := ctx.People.Create(p)

	if err != nil {
		t.Error(err)
	}

	err = ctx.People.Delete(record.GetKey())

	if err != nil {
		t.Error(err)
	}

	_, rerr := ctx.People.FindByKey(record.GetKey())

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
	rec, _ := ctx.People.Create(p1)
	ctx.People.Create(p2)

	result, err := ctx.People.FindFirst(func(obj husk.Dataer) bool {
		return obj.(*sample.Person).Name == "Sarel"
	})

	if err != nil {
		t.Error(err)
	}

	firstID := result.GetKey()

	if firstID != rec.GetKey() {
		t.Errorf("Wrong ID, Expected %v, got %v", rec.GetKey(), firstID)
	}
}

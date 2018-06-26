package tests

import (
	"fmt"
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
	err := husk.DestroyContents("./db")

	if err != nil {
		log.Println(err)
	}
}

func TestCreate_MustPersist(t *testing.T) {
	//defer DestroyData()

	p := sample.Person{"Jan", 25}

	record, err := ctx.People.Create(&p)

	if err != nil {
		t.Error(err)
	}

	againP, ferr := ctx.People.FindByID(record.GetID())

	if ferr != nil {
		t.Error(ferr)
	}

	if againP == nil {
		t.Error("Record not found")
	}
}

func TestCreate_MultipleEntries_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{"Johan", 13}
	p1 := sample.Person{"Sarel", 15}
	p2 := sample.Person{"Jaco", 24}

	ctx.People.Create(p)
	ctx.People.Create(p1)
	p2Record, rerr := ctx.People.Create(p2)

	if rerr != nil {
		t.Error(rerr)
	}

	_, err := ctx.People.FindByID(p2Record.GetID())

	if err != nil {
		t.Error(err)
	}
}

func TestUpdate_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{"Sarie", 45}

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

	againP, ferr := ctx.People.FindByID(record.GetID())

	if ferr != nil {
		t.Error(ferr)
	}

	againData := againP.Data().(*sample.Person)

	if againData.Age != p.Age {
		t.Errorf("Expected %v, got %v", p.Age, againData.Age)
	}
}

func TestDelete_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{"DeleteMe", 67}

	record, err := ctx.People.Create(p)

	if err != nil {
		t.Error(err)
	}

	err = ctx.People.Delete(record.GetID())

	if err != nil {
		t.Error(err)
	}

	_, rerr := ctx.People.FindByID(record.GetID())
	expectedErr := fmt.Sprintf("ID %v not found in table Person", record.GetID())

	if rerr.Error() != expectedErr {
		t.Error("Expected item to be deleted.")
	}
}

func TestFind_FindFilteredItems(t *testing.T) {
	defer DestroyData()

	p := sample.Person{"Johan", 13}
	p1 := sample.Person{"Sarel", 15}
	p2 := sample.Person{"Jaco", 24}

	ctx.People.Create(p)
	rec, _ := ctx.People.Create(p1)
	ctx.People.Create(p2)

	results := ctx.People.Find(1, 1, func(obj husk.Dataer) bool {
		return obj.(*sample.Person).Name == "Sarel"
	})

	if len(results) != 1 {
		t.Errorf("Expected %v, got %v", 1, len(results))
	}

	firstID := results[0].GetID()

	if firstID != rec.GetID() {
		t.Errorf("Wrong ID, Expected %v, got %v", rec.GetID(), firstID)
	}
}

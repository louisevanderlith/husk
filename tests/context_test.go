package tests

import (
	"encoding/json"
	"log"
	"strings"
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

	if set.Error != nil {
		t.Error(set.Error)
	}

	ctx.People.Save()
	recKey := set.Record.GetKey().String()
	k, _ := husk.ParseKey(recKey)
	againP, ferr := ctx.People.FindByKey(k)

	if ferr != nil {
		t.Error(ferr)
		return
	}

	if againP == nil {
		t.Error("Record not found")
		return
	}

	if againP.GetKey() != k {
		t.Errorf("Expected %s, %s", recKey, againP.GetKey())
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
	ctx.People.Save()

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
	ctx.People.Save()

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
	defer ctx.People.Save()

	if set.Error != nil {
		t.Error(set.Record)
	}

	firstUpdate := set.Record.Meta().LastUpdated()

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

	if againMeta.LastUpdated() == firstUpdate {
		t.Errorf("Expected %v, got %v", firstUpdate, againMeta.LastUpdated())
	}
}

func TestDelete_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "DeleteMe", Age: 67}
	set := ctx.People.Create(p)

	if set.Error != nil {
		t.Error(set)
		return
	}

	ctx.People.Save()
	t.Log(set.Record.GetKey())
	//
	err := ctx.People.Delete(set.Record.GetKey())

	if err != nil {
		t.Error(err)
		return
	}

	_, rerr := ctx.People.FindByKey(set.Record.GetKey())

	if rerr == nil {
		t.Error("Expected item to be deleted. 'Not found error...'")
	}
}

func TestFind_FindFilteredItems(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Johan", Age: 13}
	p1 := &sample.Person{Name: "Sarel", Age: 15}
	p2 := sample.Person{Name: "Jaco", Age: 24}

	ctx.People.Create(&p)
	set := ctx.People.Create(p1)
	ctx.People.Create(&p2)
	ctx.People.Save()

	ctx.People.Update(set.Record)
	ctx.People.Save()
	result, err := ctx.People.FindFirst(sample.ByName("Sarel"))
	//result := ctx.People.Find(1, 3, husk.Everything())

	if err != nil {
		t.Error(err)
		return
	}

	if result.GetKey() != set.Record.GetKey() {
		t.Errorf("Wrong ID, Expected %v, got %v", set.Record.GetKey(), result.GetKey())
	}
}

func TestCalc_SumTotalBalance(t *testing.T) {
	bal := float32(0)

	ctx.People.Calculate(&bal, sample.SumBalance())

	t.Log(bal)
	t.Fail()
	if bal == 0 {
		t.Error("balance not updated")
	}
}

func TestCalc_FindLowestBalance(t *testing.T) {
	name := ""

	ctx.People.Calculate(&name, sample.LowestBalance())

	if name != "Kelley" {
		t.Errorf("name not found, got %s", name)
	}
}

func TestExsits_Empty_MustTrue(t *testing.T) {
	defer DestroyData()

	expect := true
	actual := ctx.People.Exists(husk.Everything())

	if actual != expect {
		t.Errorf("Expecting %v, got %v", expect, actual)
	}
}

func TestExsits_Any_MustTrue(t *testing.T) {
	defer DestroyData()
	p := sample.Person{Name: "Weirdo", Age: 55}
	ctx.People.Create(p)

	expect := true
	actual := ctx.People.Exists(husk.Everything())

	if actual != expect {
		t.Errorf("Expecting %v, got %v", expect, actual)
	}
}

func TestFilter_FindWarden_MustBe10(t *testing.T) {
	records := ctx.People.Find(1, 11, sample.ByNameAndAge("Warden", 22))

	if records.Count() != 10 {
		t.Errorf("Expected %d records, got %d", 10, records.Count())
	}
}

func TestFilter_FindEverything_MustBe1000(t *testing.T) {
	records := ctx.People.Find(1, 1000, husk.Everything())

	if records.Count() != 1000 {
		t.Errorf("Expecting 1000, got %d", records.Count())
	}
}

func TestRecordSet_ToJSON_MustBeClean(t *testing.T) {
	rows := ctx.People.Find(1, 5, husk.Everything())
	bits, _ := json.Marshal(rows)

	if strings.Contains(string(bits), "Value") {
		t.Error("Final Object has Value")
	}
}

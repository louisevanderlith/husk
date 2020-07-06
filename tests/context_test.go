package tests

import (
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/tests/sample"
	"log"
	"testing"
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

func TestDiscover_ListNames(t *testing.T) {
	defer DestroyData()
	ctx.Seed()

	exp := []string{"Journals"}
	act := husk.TableNames(ctx)

	if len(act) != len(exp) {
		t.Error("invalid length discovered")
		return
	}

	if exp[0] != act[0] {
		t.Errorf("Expected %s, found %s", exp, act)
	}
}

func TestDiscover_ListLayouts(t *testing.T) {
	defer DestroyData()
	act := husk.TableLayouts(ctx)

	if len(act) != 1 {
		t.Error("invalid length discovered")
		return
	}

	if act["Journals"] == nil {
		t.Errorf("no object found %v", act)
	}
}

func TestFind_SearchItems(t *testing.T) {
	set, err := benchCtx.Journals.Find(1, 10, sample.ByPublisher("Universidade Federal do Rio Grande"))

	if err != nil {
		t.Error(err)
		return
	}

	itor := set.GetEnumerator()

	for itor.MoveNext() {
		curr := itor.Current()

		t.Log(curr.Data())
	}

	if set.Count() != 6 {
		t.Errorf("%+v\n", set.Count())
	}
}

/*
func TestCreate_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Journal{Name: "Jan", Age: 25}

	set := ctx.Journals.Create(p)

	if set.Error != nil {
		t.Error(set.Error)
	}

	err := ctx.Journals.Save()

	if err != nil {
		t.Error("Save Error", err)
		return
	}

	recKey := set.Record.GetKey().String()
	k, _ := husk.ParseKey(recKey)
	againP, err := ctx.Journals.FindByKey(k)

	if err != nil {
		t.Error("Find Error", err)
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

	ctx.Journals.Create(p)
	ctx.Journals.Create(p1)
	p2Set := ctx.Journals.Create(p2)
	ctx.Journals.Save()

	if p2Set.Error != nil {
		t.Error(p2Set.Error)
	}

	_, err := ctx.Journals.FindByKey(p2Set.Record.GetKey())

	if err != nil {
		t.Error(err)
	}
}

func TestUpdate_MustPersist(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Sarie", Age: 45}

	set := ctx.Journals.Create(p)

	if set.Error != nil {
		t.Error(set)
		return
	}

	err := ctx.Journals.Save()

	if err != nil {
		t.Error(err)
	}

	pData := set.Record.Data().(sample.Person)
	pData.Age = 67

	err = ctx.Journals.Update(set.Record)

	if err != nil {
		t.Error(err)
		return
	}

	againP, ferr := ctx.Journals.FindByKey(set.Record.GetKey())

	if ferr != nil {
		t.Error(ferr)
	}

	againData := againP.Data().(sample.Person)

	if againData.Age != p.Age {
		t.Errorf("Expected %v, got %v", p.Age, againData.Age)
	}
}

func TestUpdate_LastUpdatedMustChange(t *testing.T) {
	defer DestroyData()

	p := sample.Person{Name: "Sarie", Age: 45}

	set := ctx.Journals.Create(p)

	if set.Error != nil {
		t.Error(set)
		return
	}

	err := ctx.Journals.Save()

	if err != nil {
		t.Error(err)
		return
	}

	firstUpdate := set.Record.Meta().LastUpdated()

	pData := set.Record.Data().(sample.Person)
	pData.Age = 67

	err = ctx.Journals.Update(set.Record)

	if err != nil {
		t.Error(err)
		return
	}

	againP, ferr := ctx.Journals.FindByKey(set.Record.GetKey())

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
	set := ctx.Journals.Create(p)

	if set.Error != nil {
		t.Error(set)
		return
	}

	ctx.Journals.Save()
	t.Log(set.Record.GetKey())
	//
	err := ctx.Journals.Delete(set.Record.GetKey())

	if err != nil {
		t.Error(err)
		return
	}

	_, rerr := ctx.Journals.FindByKey(set.Record.GetKey())

	if rerr == nil {
		t.Error("Expected item to be deleted. 'Not found error...'")
	}
}

func TestFind_FindFilteredItems(t *testing.T) {
	//defer DestroyData()

	p := sample.Person{Name: "Johan", Age: 13}
	p1 := sample.Person{Name: "Sarel", Age: 15}
	p2 := sample.Person{Name: "Jaco", Age: 24}

	ctx.Journals.Create(p)
	set := ctx.Journals.Create(p1)

	if set.Error != nil {
		t.Fatal(set.Error)
	}

	ctx.Journals.Create(p2)
	ctx.Journals.Save()

	err := ctx.Journals.Update(set.Record)

	if err != nil {
		t.Fatal(err)
	}

	err = ctx.Journals.Save()

	if err != nil {
		t.Fatal(err)
	}

	result, err := ctx.Journals.Find(1, 3, sample.ByName("Sarel"))

	if err != nil {
		t.Error(err)
		return
	}

	if !result.Any() {
		t.Error("no results")
	}

	itor := result.GetEnumerator()
	matchFound := false
	for itor.MoveNext() {
		if itor.Current().GetKey() == set.Record.GetKey() {
			matchFound = true
			break
		}
	}

	if !matchFound {
		t.Error("no matches found")
	}
}

func TestCalc_SumTotalBalance(t *testing.T) {
	bal := float32(0)

	ctx.Journals.Calculate(&bal, sample.SumBalance())

	t.Log(bal)
	if bal == 0 {
		t.Error("balance not updated")
	}
}

func TestUsers_FindEverything(t *testing.T) {
	rset, err := ctx.Journals.Find(1, 10, husk.Everything())

	if err != nil {
		t.Fatal(err)
		return
	}

	if !rset.Any() {
		t.Error("no data found")
	}
}

func TestCalc_FindLowestBalance(t *testing.T) {
	name := ""

	ctx.Journals.Calculate(&name, sample.LowestBalance())

	if name != "Kelley" {
		t.Errorf("name Kelley not found, got %s", name)
	}
}

func TestExsits_Empty_MustTrue(t *testing.T) {
	defer DestroyData()

	expect := true
	actual := ctx.Journals.Exists(husk.Everything())

	if actual != expect {
		t.Errorf("Expecting %v, got %v", expect, actual)
	}
}

func TestExsits_Any_MustTrue(t *testing.T) {
	//defer DestroyData()
	p := sample.Person{Name: "Weirdo", Age: 55}
	ctx.Journals.Create(p)

	expect := true
	actual := ctx.People.Exists(husk.Everything())

	if actual != expect {
		t.Errorf("Expecting %v, got %v", expect, actual)
	}
}

func TestFilter_FindWarden_MustBe10(t *testing.T) {
	parm := sample.Person{
		Name: "Warden",
		Age:  22,
	}

	records, err := ctx.Journals.Find(1, 11, sample.ByObject(parm))

	if err != nil {
		t.Fatal(err)
		return
	}

	if records.Count() != 10 {
		t.Errorf("Expected %d records, got %d", 10, records.Count())
	}
}

func TestFilter_FindWarden_MustBeByFields(t *testing.T) {
	parm := sample.Person{
		Name: "Warden",
		Age:  22,
	}

	records, err := ctx.Journals.Find(1, 10, husk.ByFields(parm))

	if err != nil {
		t.Fatal(err)
		return
	}

	enumer := records.GetEnumerator()
	for enumer.MoveNext() {
		t.Log(enumer.Current())
	}

	if records.Count() != 10 {
		t.Errorf("Expected %d records, got %d", 10, records.Count())
	}
}

func TestFilter_FindEverything_MustBe1000(t *testing.T) {
	records, err := ctx.Journals.Find(1, 1000, husk.Everything())

	if err != nil {
		t.Fatal(err)
		return
	}

	if records.Count() != 1000 {
		t.Errorf("Expecting 1000, got %d", records.Count())
	}
}

func TestRecordSet_ToJSON_MustBeClean(t *testing.T) {
	rows, err := ctx.Journals.Find(1, 5, husk.Everything())

	if err != nil {
		t.Fatal(err)
		return
	}

	bits, _ := json.Marshal(rows)

	if strings.Contains(string(bits), "Value") {
		t.Error("Final Object has Value")
	}
}
*/

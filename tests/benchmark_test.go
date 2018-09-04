package tests

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/tests/sample"
)

var (
	benchCtx  sample.Context
	bodies    []sample.Person
	idx       int
	bodycount int
)

func init() {
	benchCtx = sample.NewContext()
	loadBodies()
	idx = 0
	bodycount = len(bodies)
}

func loadBodies() {
	var result []sample.Person
	byts, err := ioutil.ReadFile("raw_data.json")

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(byts, &result)

	if err != nil {
		panic(err)
	}

	bodies = result
}

func TestInserts_SampleETL(t *testing.T) {

	for nxtPerson, valid := getNextPerson(); valid; nxtPerson, valid = getNextPerson() {

		_, err := benchCtx.People.Create(nxtPerson)

		if err != nil {
			t.Error(err)
		}
	}

	//t.Fail()

}

// BenchmarkInserts run a benchmark with simple objects
func BenchmarkInserts(b *testing.B) {

	nxtPerson, valid := getNextPerson()

	if valid {

		benchCtx.People.Create(nxtPerson)
	}
}

func BenchmarkFilter_HighBalance(b *testing.B) {
	set := benchCtx.People.Find(1, 50, func(data husk.Dataer) bool {
		obj := data.(*sample.Person)
		b.Logf("%+v\n", obj)
		for _, v := range obj.Accounts {
			if v.Balance > 30000 {
				return true
			}
		}

		return false
	})

	if set.Count() < 5 {
		b.Errorf("%+v\n", *set)
	}
}

func getNextPerson() (*sample.Person, bool) {
	if idx == bodycount {
		return nil, false
	}

	result := bodies[idx]
	fmt.Println("IDX:", idx)
	idx++

	return &result, true
}

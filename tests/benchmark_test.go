package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/tests/sample"
)

var (
	benchCtx  sample.Context
	bodies    []sample.Person
	idx       int
	bodycount int
	counter   int
	startTime time.Time
)

func init() {
	benchCtx = sample.NewContext()
	loadBodies()
	idx = 0
	bodycount = len(bodies)
	counter = 0
	startTime = time.Now()
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
	defer benchCtx.People.Save()
	count := 0

	nxtPerson, _ := getNextPerson()

	for time.Since(startTime) < time.Second*20 {
		count++
		set := benchCtx.People.Create(nxtPerson)

		if set.Error != nil {
			t.Error(set.Error)
		}
	}

	t.Logf("Completed: %d", count)
	t.Fail()
}

// BenchmarkInserts run a benchmark with simple objects
func BenchmarkInserts(b *testing.B) {

	nxtPerson, valid := getNextPerson()

	if valid {
		benchCtx.People.Create(nxtPerson)
	}

	defer benchCtx.People.Save()
}

func BenchmarkFilter_HighBalance(b *testing.B) {
	set := benchCtx.People.Find(1, 50, func(data husk.Dataer) bool {
		obj := data.(*sample.Person)

		for _, v := range obj.Accounts {
			if v.Balance > 50 {
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
	idx++

	return &result, true
}

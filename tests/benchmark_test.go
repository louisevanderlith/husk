package tests

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
	"time"

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

	log.Printf("Loaded %d records", len(result))
	bodies = result
}

//Warning: Inserts the same record for 20seconds.
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

func TestInserts_SampleData_MustLoadAll(t *testing.T) {
	defer benchCtx.People.Save()

	count := 0

	for nxtPerson, valid := getNextPerson(); valid; nxtPerson, valid = getNextPerson() {
		count++

		set := benchCtx.People.Create(nxtPerson)

		if set.Error != nil {
			t.Error(set.Error)
		}
	}

	log.Printf("Completed %d\n", count)
}

// BenchmarkInserts run a benchmark with simple objects
func BenchmarkInserts(b *testing.B) {

	nxtPerson, valid := getNextPerson()

	if valid {
		set := benchCtx.People.Create(nxtPerson)

		if set.Error != nil {
			b.Logf("Record Failed to create: %s", set.Error.Error())
		}
	}

	defer benchCtx.People.Save()
}

func BenchmarkFilter_HighBalance(b *testing.B) {
	set := benchCtx.People.Find(1, 50, sample.HigherBalance(123.50))

	if set.Count() < 5 {
		b.Errorf("%+v\n", set)
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

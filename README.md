# Husk DB-Engine
Husk was designed to be used directly by Webservice and API based applications.
The amount of traffic would then be limited by the amount of requests your API can handle, instead of how-many the Database server can take.

Thie engine attempts to force users to keep business logic close to their required objects, and minimizes entry points and chances for loop holes. ie. when a developer accesses or modifies data outside of the intended scope. 

This includes many sources; 
* Directly modifying data via an external tool. (SSMS, WorkBench, etc.)
* Able to access core "Database" API in higher-level code, rather than the Logic Layer as intended.  (Front-end, Public facing API End-points)

All records are internally sorted by their creation Timestamp, and then their traditional "ID".
This combination is refered to as a Key.

The database engine works similiar to ISAM, as it stores data on a sequential "tape" which lives in memory and on disk.
Husk uses an index with pointers to the actual location on the tape for faster access.

# Bench History (TestInserts_SampleETL):
Please note these numbers come from our Sample_ETL test, which inserts the same record(16kb) for 20seconds
* "0.1 (One Record, One File) Write: 138rec/s"
* 0.2 (BigFile) Write: 509rec/s (x3.6)
* 0.3 (Dump Index only on save) Write: 1463rec/s (x3)
* 0.4 (Better File handling) Write: 1221rec/s (0%)
* 0.5 (Index Refactor, keys are Ptrs, improved read) Write: 2315rec/s (x2)
* 0.6 (Key isn't a pointer anymore) 4314rec/s

# Average performance
* MAC 3167rec/s (Unicorn Power)
* WINDOWS 2315/rec/s (Spinning Disk, AMD)
* LINUX 2289rec/s (SSD, Intel i5(2nd gen))

# Database Engine
* Data-orientation and clustering
* Everything related to an object will always remain nested within that object. 
* ISAM?

# Setting up a database
Create a Table Object
```go 
package sample

import "github.com/louisevanderlith/husk"

//Person Data Record
type Person struct {
	Name     string `hsk:"size(50)"`
	Age      int
	Accounts []Account
}

//Valid - To qualify as a Data Record, a struct MUST have a Valid function
func (o Person) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}
```

Create a context for quick access to Tables
```go
package sample

import "github.com/louisevanderlith/husk"

//Context holds the Tables we have access to 
type Context struct {
	//People table 
	People husk.Tabler
}

func NewContext() Context {
	result := Context{}

	//Creats a new "Person" Table
	result.People = husk.NewTable(new(Person))

	return result
}

func (ctx Context) Seed() {
	//Seed files can be specified, so that we have data to boot.
	err := ctx.People.Seed("people.seed.json")

	if err != nil {
		panic(err)
	}

	ctx.People.Save()
}
```

## Using the table
Create a record
```go
p := sample.Person{Name: "Jan", Age: 25}

//Send a Ptr to the object to Create
set := ctx.People.Create(&p)

if set.Error != nil {
	t.Error(set.Error)
}

//Persist the changes
ctx.People.Save()
```

Find and update
```go
//Find by it's Key
person, _ := ctx.People.FindByKey(key)
person.Age = 87

ctx.People.Update(person)

//Persist the changes
ctx.People.Save()
```
Working with collections
```go
//Find 'Everything', but I only want the first 3
result := ctx.People.Find(1, 3, husk.Everything())

if result == nil {
	return
}

//Gets the iterable collection
rator := result.GetEnumerator()

//Moves to the next item in the collection, until there isn't anything else
for rator.MoveNext() {
	curr := rator.Current()
	someone := curr.Data().(*sample.Person)

	log.Printf("$v\n", someone)
}
```

Creating filters for records
```go
//Specify a Data Filter for the given Record
type personFilter func(obj *Person) bool

//Filter is called by Husk, but casted to the correct type.
func (f personFilter) Filter(obj husk.Dataer) bool {
	return f(obj.(*Person))
}

//Filter People by their Name
func ByName(name string) personFilter {
	return func(obj *Person) bool {
		return obj.Name == name
	}
}

//Filter for searching by Balance on Accounts
func SameBalance(balance float32) personFilter {
	return func(obj *Person) bool {
		for _, v := range obj.Accounts {
			if v.Balance == balance {
				return true
			}
		}

		return false
	}
}
```
# Husk Database Engine [![CodeFactor](https://www.codefactor.io/repository/github/louisevanderlith/husk/badge)](https://www.codefactor.io/repository/github/louisevanderlith/husk)
Husk was designed to be used directly by Webservice and API based applications without the need for an external storage provider.

## Table of Contents
* What is HuskDB?
* Setup and Installation
* Creating Table/Structure
* Preparing Seed Data
* Operations
    * Create
    * Read
    * Update
    * Delete
    * Calculate
* Working with Collections
* Advantages over traditional SQL
* hsk Tags
* Benchmarks

### What is HuskDB
Husk database is an embedded, object-oriented data store which uses Go to interact with records.
It is meant to be used by small webservices which fully control the data it hosts.

This engine attempts to force users to keep business logic close to their required objects, and minimizes entry points and chances for loop holes. ie. when a developer accesses or modifies data outside of the intended scope. 

This includes many sources; 
* Directly modifying data via an external tool. (SQL Server Management Studio, WorkBench, Toad, etc.)
* Able to access core "Database" API in higher-level code, rather than the Logic Layer as intended.  (Front-end, Public facing API End-points)
* Editing the data file directly.

Creation timestamps sort records internally, and thereafter their traditional "ID".
The 'Key' is created from the combination of timestamp and id "0`0".
The Key allows the index to sort the records by creation date, by default. 
This creates faster access to the most recent records and removes need for sorting after every query.

The database engine works similar to ISAM, as it stores data on a sequential "tape" which lives on disk and held in memory.
Husk uses an index with pointers to the actual location on the tape for faster access.

### Setup and Installation
As HuskDB is a library, which can simply be imported into any Go project.

To download the Husk module, run;

```$ go get -u github.com/louisevanderlith/husk```

After installing the latest module, all that is required is to create a Table/Structure for your data. The following 
section will provide information on setting up the Table.

More usage examples can be found under the /tests folder.

### Creating Table/Structure
* Define data object
Any type which has the 'Valid()' function, qualifies as a data object, and can be used as a Tabler.

```go
package sample

import "github.com/louisevanderlith/husk"

//Person Data Record
type Person struct {
	Name     string `hsk:"size(50)"`
	Age      int
	Accounts []Account
}

//Valid - To qualify as a Data Record, 
//a struct MUST have a Valid function
func (o Person) Valid() (bool, error) {
	return husk.ValidateStruct(&o)
}
```
In this sample, we create a "Person" object which holds Name, Age and Account information.

The Name property can only have a value which has a maximum length of 50characters. More information in `hsk` tags can be found later in this document.

* Create a context for the structure
```go
    package sample
    
    import (
        "github.com/louisevanderlith/husk"
        "github.com/louisevanderlith/husk/serials"
    )
    
    //Context holds the Tables we want to access
    type Context struct {
        //People table 
        People husk.Tabler
    }
    
    // NewContext returns the context object with Table values initialized
    func NewContext() Context {
        result := Context{}
    
        //Creats a new "Person" Table, with GobSerializer
        result.People = husk.NewTable(Person{}, serials.GobSerial{})
    
        return result
    }
```

We now have a context which can be utilised by our application

### Preparing Seed Data
Husk is able to import seed data via a JSON file.
* Create seed file (people.seed.json) in a folder named 'db'
* Populate seed data.
    Seed data should be provided as a JSON array. 
    Example;
    ```json
    [
        {
            "Name": "Charlie",
            "Age": 23,
            "Accounts": [{"Number": 1234500,"Balance": 0.10, "Transactions": []}]
        },
        {
            "Name": "Mike",
            "Age": 48,
            "Accounts": [{"Number": 1234501,"Balance": 99.10, "Transactions": []}]
        }
    ]
  ```
* Call the Seed function on structures to populate the table

```go 
    func (ctx Context) Seed() {
        //Seed files can be specified, so that we have data to boot.
        err := ctx.People.Seed("people.seed.json")
    
        if err != nil {
            panic(err)
        }
    
        ctx.People.Save()
    }
```

### Operations
* Create
Records can be created by providing an object to the Create function.
Husk also supports .CreateMulti which can be used to create many records at once.
```go
    p := sample.Person{Name: "Jimmy", Age: 25}
    
    //Send the object to the context for creation
    set := ctx.People.Create(p)
    
    if set.Error != nil {
        t.Error(set.Error)
    }
    
    //Persist the changes
    ctx.People.Save()
```

* Read

Records can be located using various methods:
    
1. By Key
    
This is the fastest way to locate a record
    
```go
    //parse the string representation of the key to an actual husk.Key
    k, err := husk.ParseKey("0`0")
    
    //find the record with the key
    rec, err := ctx.People.FindByKey(key)
```
    
2. By Filter
    
Filters can be defined as functions which return true/false depending on the conditions.
Husk also provides a default filter "husk.Everything" which can be used to return all records.
    
```go
    type personFilter func(obj Person) bool
    
    func (f personFilter) Filter(obj husk.Dataer) bool {
        return f(obj.(Person))
    }
    
    func ByName(name string) personFilter {
        return func(obj Person) bool {
            return obj.Name == name
        }
    }
```

Operations that find records always have to supply page size and index.
FindFirst is a shortcut for Find(1,1, filter)
```go
    //Find People 'ByName', but I only want the first 3 matches
    result := ctx.People.Find(1, 3, ByName("Jimmy"))
    
    //Find 'All' People, but I only want the first 3 matches
    result = ctx.People.Find(1, 3, husk.Everything())
```

* Update

```go
    //First find the desired record, in this case "ByKey"
    p, _ := ctx.People.FindByKey(key)
    p.Age = 87
    
    ctx.People.Update(p)
    
    //Persist the changes
    ctx.People.Save()
```
* Delete

Records can be deleted directly from the database by providing the key of the record to be removed.
```go
    //Provide the key to the delete function
    err := ctx.People.Delete(key)
    
    //Persist the changes
    ctx.People.Save()
```

* Calculate

This function can be used in multiple ways to generate datasets.
There is no concept of "SELECT" in Husk, but it still provides a way of creating custom datasets. 
When thinking in terms of traditional SQL, this would be a Stored Procedure.
The calculate function could be used to generate custom datasets.  

Calculators can be defined in the same way as Filters, as they function in exactly the same manner. 
For this reason, we can chain Filters and Calculators to narrow result sets.

```go
    type personCalc func(result interface{}, obj Person) error
    
    // Calc can take a pointer to a result, and update it with values found in the data obj
    func (f personCalc) Calc(result interface{}, obj husk.Dataer) error {
        return f(result, obj.(Person))
    }
``` 

In the following example we want to get the Total value of a Person's Accounts.

```go
// SumBalance will return the SUM Balance of a Person's Accounts
func SumBalance() personCalc {
	return func(result interface{}, obj Person) error {
		answ := float32(0)
		for _, acc := range obj.Accounts {
			answ += acc.Balance
		}

		totl := result.(*float32)
		*totl += answ

		return nil
	}
}
```

This example shows how we can get the Lowest Valued Account
```go
func LowestBalance() personCalc {
	min := float32(9999999)
	return func(result interface{}, obj Person) error {
        // We can utilise previously defined functions
        balnce := float32(0)
		err := SumBalance(&balnce)
        
		if answ < min {
			min = answ
			n := result.(*string)
			*n = obj.Name
		}

		return nil
	}
}
```

### Working with Collections
After finding records, you may want to iterate over them to perform other operations.
Husk provides a way of enumerating these collections.
```go
//Find 'Everything', but I only want the first 3
result := ctx.People.Find(1, 3, husk.Everything())

//Gets the iterable collection
rator := result.GetEnumerator()

//Moves to the next item in the collection, until there isn't anything else
for rator.MoveNext() {
	curr := rator.Current()
	someone := curr.Data().(sample.Person)

	log.Printf("$v\n", someone)
}
``` 

### Advantages over traditional SQL
Husk doesn't use any Query language, and all data manipulation happens via the code. 
The immediately avoids SQL injection attacks, and narrows the points of entry.

* No SELECT, Datasets are 'Calculated' and relates more to Stored Procedures, but written in Go.
* No WHERE Filter using Go functions. Functions can be easily tested and compile with the application.
* Index Keys are separate from Data, and should only be used to Identify records and create relationships across micro-services.   
* The database stores records in descending order of creation.  
* The Primary Key consists of a Timestamp and ID. Quite husk.Key
* Querying collections always require the 'pagesize' (current page and results per page) to be specified.
    This forces the developer to always limit the amount of data returned, thus reducing query times.
* Database embedded into application, no externally hosted services. Greatly reduces the attack surface of the data. 
* Doesn't require a "ConnectionString"
* TDD and Unit Tests can use Husk to create tables that work with the logic. Object-Oriented approach to working with data.
* Tables can be mapped directly to JSON structures. Big JSON files can be easily analysed.
* Seed files are JSON documents, so other systems can easily be migrated.  
* Objects are "Serialization" aware, and columns like 'Password' can easily be hidden using `json:"-"`
* Everything related to an object will always remain nested within that object.
* Database will only consist of one or two tables, micro-services should only control one part of a system.

### Benchmarks
History:
Please note these numbers come from our Sample_ETL test, which inserts the same record(16kb) for 20seconds 
(This function has since been deprecated, and a better benchmark is yet to be written)
* v1.0.1 Write: 138rec/s
Every record saved creates a new file. Very slow read, and write.

* v1.0.2 Write: 509rec/s (x3.6 improvement)
One file stores all records. Greatly improves reads. 

* v1.0.3 Write: 1463rec/s (x3 improvement)
Index file will only be updated on disk when the context gets saved.

* v1.0.4 Write: 1221rec/s (x1 improvement)
File operations improved. 

* v1.0.5 Write: 2315rec/s (x2 improvement)
Indexing logic refactored and Keys changed to Pointers. Improved reading.

* v1.0.6 Write: 4314rec/s (x2 improvement)
Keys are no longer Pointers.

* v1.0.7 Write: Unknown
General improvements

Average Write Performance:

* MAC 3167rec/s (Unicorn Power)
* WINDOWS 2315/rec/s (Spinning Disk, AMD)
* LINUX 2289rec/s (SSD, Intel i5(2nd gen))


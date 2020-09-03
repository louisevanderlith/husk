package hsk

import (
	"github.com/louisevanderlith/husk/validation"
	"io"
)

//Table is used to interact with records
type Table interface {
	Name() string
	//Exists confirms the existence of a record
	Exists(filter Filter) bool

	//FindByKey finds a record with a matching key.
	FindByKey(key Key) (Record, error)
	//Find looks for records that match the filter.
	Find(page, pageSize int, filter Filter) (Page, error)
	//FindFirst does what Find does, but will only return one record.
	FindFirst(filter Filter) (Record, error)

	//Map can modify a result set with data values
	Map(result interface{}, calculator Mapper) error

	//Create saves a new object to the database
	Create(obj validation.Dataer) (Key, error)
	//CreateMulti saves multiple records
	CreateMulti(objs ...validation.Dataer) ([]Key, error)
	//Update records changes made to a record.
	Update(key Key, obj validation.Dataer) error
	//Delete removes a record with the matching key.
	Delete(key Key) error

	//Seeds data from a json file
	Seed(seedfile string) error
	//Seeds data from a io.reader
	//SeedReader(r io.Reader) error --soon

	SaveWriter(w io.Writer) error
	Save() error
}

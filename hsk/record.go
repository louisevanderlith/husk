package hsk

import (
	"encoding/json"
	"github.com/louisevanderlith/husk/validation"
)

//Record is what defines a record, and what it can do
type Record interface {
	GetKey() Key
	Data() validation.Dataer
}

//Record exists of the actual data, and meta info about the data
type record struct {
	key   Key
	value validation.Dataer
}

//MakeRecord creates a new Record
func MakeRecord(k Key, obj validation.Dataer) Record {
	return &record{k, obj}
}

//GetKey returns the key value from meta
func (r record) GetKey() Key {
	return r.key
}

//Data returns the record's actual data
func (r record) Data() validation.Dataer {
	return r.value
}

//MarshalJSON returns Records as {K:[KEY](1540921456-18), V: [VALUE](obj{})}
func (r *record) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		K string
		V interface{}
	}{r.GetKey().String(), r.value})
}

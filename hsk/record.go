package hsk

import (
	"encoding/json"
)

//Record exists of the actual data, and meta info about the data
type Record struct {
	meta  Meta
	value Dataer
}

//MakeRecord creates a new Record
func MakeRecord(meta Meta, obj Dataer) Recorder {
	return &Record{meta, obj}
}

//GetKey returns the key value from meta
func (r Record) GetKey() Key {
	return r.meta.GetKey()
}

//Meta returns the record's meta information
func (r Record) Meta() Meta {
	return r.meta
}

//Data returns the record's actual data
func (r Record) Data() Dataer {
	return r.value
}

//Set applies the value of 'obj' to the current record.
func (r *Record) Set(obj Dataer) error {
	err := obj.Valid()

	if err != nil {
		return err
	}

	r.value = obj

	return nil
}

//MarshalJSON returns Records as {K:[KEY](1540921456-18), V: [VALUE](obj{})}
func (r *Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		K string
		V interface{}
	}{r.GetKey().String(), r.value})
}

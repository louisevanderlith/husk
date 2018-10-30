package husk

import (
	"encoding/json"
)

type Record struct {
	meta  *meta
	value Dataer
}

func MakeRecord(meta *meta, obj Dataer) *Record {
	return &Record{meta, obj}
}

func (r Record) GetKey() *Key {
	return r.meta.GetKey()
}

func (r Record) Meta() *meta {
	return r.meta
}

func (r Record) Data() Dataer {
	return r.value
}

func (r *Record) Set(obj Dataer) error {
	valid, err := obj.Valid()

	if err != nil || !valid {
		return err
	}

	r.value = obj

	return nil
}

func (r *Record) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.value)
}

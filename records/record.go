package records

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/validation"
)

//Record exists of the actual data, and meta info about the data
type record struct {
	Key   hsk.Key
	Value validation.Dataer
}

//MakeRecord creates a new Record
func MakeRecord(k hsk.Key, obj validation.Dataer) hsk.Record {
	return &record{k, obj}
}

func NewRecord(t validation.Dataer) hsk.Record {
	return &record{keys.CrazyKey(), t}
}

//GetKey returns the key value from meta
func (r *record) GetKey() hsk.Key {
	return r.Key
}

//Data returns the record's actual data
func (r *record) GetValue() validation.Dataer {
	return r.Value
}

func (r record) Compare(obj interface{}) int8 {
	r2 := obj.(*record)
	return r.Key.Compare(r2.Key)
}

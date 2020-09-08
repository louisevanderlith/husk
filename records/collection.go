package records

import (
	"encoding/json"
	"errors"
	"github.com/louisevanderlith/husk/collections"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/validation"
	"reflect"
)

//Collection is a set of Records
type Collection interface {
	collections.Enumerable
	Get(index int) hsk.Record
	Add(rec hsk.Record) int
	Contains(rec hsk.Record) bool
	IndexOf(rec hsk.Record) int
	Insert(index int, rec hsk.Record)
	Remove(rec hsk.Record)
	RemoveAt(index int)
	Count() int
}

func NewCollection() Collection {
	return &collection{
	}
}

func CollectionOf(t validation.Dataer) Collection {
	return &collection{
		of: t,
	}
}

type SliceORecords []hsk.Record

func createCollection(obj validation.Dataer) SliceORecords {
	rec := NewRecord(obj)
	//coll := reflect.Zero(reflect.SliceOf(reflect.TypeOf(rec)))

	result := SliceORecords{rec}
	//val := reflect.ValueOf(&result)
	//val.Set(coll)

	//return result
	return result
}

type collection struct {
	items SliceORecords //[]hsk.Record
	of    validation.Dataer
}

func (c *collection) GetEnumerator() collections.Iterator {
	return collections.ReadOnlyList(reflect.ValueOf(c.items)).(collections.Iterator)
}

func (c *collection) Get(index int) hsk.Record {
	return c.items[index]
}

func (c *collection) Add(rec hsk.Record) int {
	i := c.Count()
	c.Insert(i, rec)
	return i
}

func (c *collection) Contains(rec hsk.Record) bool {
	return c.IndexOf(rec) != -1
}

func (c *collection) IndexOf(rec hsk.Record) int {
	for i := 0; i < c.Count(); i++ {
		items, comparble := c.items[i].(collections.Comparable)

		if comparble && items.Compare(rec) == 0 {
			return i
		}

		if c.items[i] == rec {
			return i
		}
	}

	return -1
}

func (c *collection) Insert(index int, rec hsk.Record) {
	if index < 0 || index > c.Count() {
		return
	}

	c.items = append(c.items, rec.(*record))
	copy(c.items[index+1:], c.items[index:])
	c.items[index] = rec.(*record)
}

func (c *collection) Remove(rec hsk.Record) {
	c.RemoveAt(c.IndexOf(rec))
}

func (c *collection) RemoveAt(index int) {
	if index >= 0 && index < c.Count() {
		for i := index; i < c.Count()-1; i++ {
			c.items[i] = c.items[i+1]
		}
	}
}

func (c *collection) Count() int {
	return len(c.items)
}

func NewData(obj validation.Dataer) (validation.Dataer, error) {
	inst := reflect.New(reflect.TypeOf(obj))
	var result validation.Dataer
	v := reflect.ValueOf(&result).Elem()

	if !v.CanSet() {
		return nil, errors.New("not settable")
	}

	v.Set(inst)

	return result, nil
}

func (c *collection) UnmarshalJSON(b []byte) error {
	var rows []json.RawMessage
	err := json.Unmarshal(b, &rows)

	if err != nil {
		return err
	}

	for _, v := range rows {
		data, err := NewData(c.of)

		if err != nil {
			return err
		}

		rec := NewRecord(data)
		err = json.Unmarshal(v, rec)
		c.items = append(c.items, rec)
	}

	return nil
	//return json.Unmarshal(b, &c.items)
}

func (c *collection) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.items)
}

package husk

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

type Table struct {
	t     reflect.Type
	name  string
	index *Index
}

func init() {
	ensureDbDirectory()
}

func NewTable(obj Dataer) Tabler {
	t := reflect.TypeOf(obj).Elem()
	name := t.Name()
	path := getIndexName(name)

	ensureTableIndex(name, path)
	index := LoadIndex(path)

	return Table{
		t:     t,
		name:  name,
		index: index,
	}
}

func (t Table) FindByKey(key Key) (Recorder, error) {
	var result Recorder
	meta := t.index.getAt(key)

	if meta == nil {
		msg := fmt.Sprintf("Key %v not found in table %s", key, t.name)

		return result, errors.New(msg)
	}

	dataObj := resultObject(t.t)
	err := read(meta.FileName, dataObj)

	if err == nil {
		result = MakeRecord(meta, dataObj)
	}

	return result, err
}

func (t Table) Find(page, pageSize int, filter Filter) *RecordSet {
	result := NewRecordSet()

	skipCount := (page - 1) * pageSize

	for t.index.MoveNext() {
		meta := t.index.Current()

		// hotfields... before file scan

		dataObj := resultObject(t.t)
		err := read(meta.FileName, dataObj)

		if err == nil && filter(dataObj) {
			if skipCount == 0 && result.Count() < pageSize {
				record := MakeRecord(meta, dataObj)
				result.Add(record)
			} else {
				skipCount--
			}
		}
	}

	return result
}

func (t Table) FindFirst(filter Filter) (Recorder, error) {
	res := t.Find(1, 1, filter)

	if !res.MoveNext() {
		return nil, errors.New("no results found")
	}

	return res.Current()
}

func (t Table) Exists(filter Filter) (bool, error) {
	item, err := t.FindFirst(filter)
	found := item != nil

	if err != nil {
		return true, err
	}

	return found, err
}

func (t Table) Create(obj Dataer) (record Recorder, err error) {
	var valid bool
	valid, err = obj.Valid()

	if valid {
		nxtKey := t.index.nextKey()

		record = NewRecord(t.name, nxtKey, obj)
		meta := record.Meta()
		err = write(meta.FileName, record.Data())

		if err == nil {
			t.index.addMeta(meta)
			t.index.dump(t.name)
		}
	}

	return record, err
}

func (t Table) Update(record Recorder) error {
	valid, err := record.Data().Valid()

	if valid {
		meta := record.Meta()
		err = write(meta.FileName, record.Data())

		if err == nil {
			meta.Updated()
			t.index.dump(t.name)
		}
	}

	return err
}

func (t Table) Delete(key Key) error {
	recMeta := t.index.getAt(key)

	if recMeta != nil {
		t.index.disable(recMeta)
		t.index.dump(t.name)
	}

	return nil
}

func resultObject(t reflect.Type) Dataer {
	return reflect.New(t).Interface().(Dataer)
}

func ensureTableIndex(tableName, indexName string) {
	created := createFile(indexName)

	if !created {
		log.Println("couldn't create index for " + tableName)
	}
}

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

func (t Table) FindByID(id int64) (Recorder, error) {
	var result Recorder
	meta := t.index.getAt(id)

	if meta == nil {
		msg := fmt.Sprintf("ID %v not found in table %s", id, t.name)

		return result, errors.New(msg)
	}

	dataObj := resultObject(t.t)
	err := read(meta.FileName, dataObj)

	if err == nil {
		result = MakeRecord(meta, dataObj)
	}

	return result, err
}

func (t Table) Find(page, pageSize int, filter Filter) Enumerator {
	var items []Recorder

	skipCount := (page - 1) * pageSize

	for _, v := range *t.index {
		dataObj := resultObject(t.t)
		err := read(v.FileName, dataObj)

		if err == nil && filter(dataObj) {
			if skipCount == 0 && len(items) < pageSize {
				record := MakeRecord(v, dataObj)
				items = append(items, record)
			} else {
				skipCount--
			}
		}
	}

	return newEnumerable(items)
}

func (t Table) FindFirst(filter Filter) (Recorder, error) {
	res := t.Find(1, 1, filter)

	return res.Next()
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
		nxtID := t.index.nextID()

		record = NewRecord(t.name, nxtID, obj)
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

func (t Table) Delete(id int64) error {
	recMeta := t.index.getAt(id)

	if recMeta != nil {
		recMeta.Disable()
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

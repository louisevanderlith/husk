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
	index Indexer
	tape  Taper
}

func init() {
	ensureDbDirectory()
}

func NewTable(obj Dataer) Tabler {
	t := reflect.TypeOf(obj).Elem()
	name := t.Name()
	path := getIndexName(name)
	trackName := getRecordName(name)

	ensureTableIndex(name, path)
	index := loadIndex(path)

	return Table{
		t:     t,
		name:  name,
		index: index,
		tape:  NewTape(trackName),
	}
}

func (t Table) FindByKey(key *Key) (Recorder, error) {
	var result Recorder
	meta := t.index.Get(key)

	if meta == nil {
		msg := fmt.Sprintf("Key %v not found in table %s", key, t.name)

		return result, errors.New(msg)
	}

	dataObj := resultObject(t.t)
	err := t.tape.Read(meta.Point(), dataObj)

	if err == nil {
		result = MakeRecord(meta, dataObj)
	}

	return result, err
}

func (t Table) Find(page, pageSize int, filter Filterer) Collection {
	result := NewRecordSet()
	skipCount := (page - 1) * pageSize

	for _, meta := range t.index.Items() {
		dataObj := resultObject(t.t)
		err := t.tape.Read(meta.Point(), dataObj)

		if err != nil {
			panic(err)
		}

		if filter.Filter(dataObj) {
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

func (t Table) FindFirst(filter Filterer) Recorder {
	res := t.Find(1, 1, filter)

	rator := res.GetEnumerator()
	if !rator.MoveNext() {
		return nil
	}

	return rator.Current()
}

func (t Table) Exists(filter Filterer) bool {
	item := t.FindFirst(filter)

	return item != nil
}

func (t Table) Create(obj Dataer) CreateSet {
	valid, err := obj.Valid()

	if !valid {
		return CreateSet{nil, err}
	}

	point, err := t.tape.Write(obj)

	if err != nil {
		return CreateSet{nil, err}
	}

	meta := t.index.CreateSpace(point)
	record := MakeRecord(meta, obj)

	t.index.Insert(meta)

	return CreateSet{record, err}
}

func (t Table) CreateMulti(objs ...Dataer) []CreateSet {
	var result []CreateSet

	for _, obj := range objs {
		set := t.Create(obj)
		result = append(result, set)
	}

	return result
}

func (t Table) Update(record Recorder) error {
	valid, err := record.Data().Valid()

	if valid {
		meta := record.Meta()
		point, err := t.tape.Write(record.Data())

		if err == nil {
			meta.Updated(point)
		}
	}

	return err
}

func (t Table) Delete(key *Key) error {
	deleted := t.index.Delete(key)

	if !deleted {
		return errors.New("nothing deleted")
	}

	return nil
}

func (t Table) Save() {
	indexName := getIndexName(t.name)

	err := write(indexName, t.index)

	if err != nil {
		panic(err)
	}
}

func resultObject(t reflect.Type) Dataer {
	return reflect.New(t).Interface().(Dataer)
}

func ensureTableIndex(tableName, indexName string) bool {
	created := createFile(indexName)

	if !created {
		log.Println("couldn't create index for " + tableName)
	}

	return created
}

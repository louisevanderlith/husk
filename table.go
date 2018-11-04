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

//NewTable returns
func NewTable(obj Dataer) Tabler {
	ensureDbDirectory()

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

func (t Table) FindByKey(key Key) (Recorder, error) {
	var result Recorder
	meta := t.index.Get(key)

	if meta == nil {
		msg := fmt.Sprintf("Key %v not found in %s", key, t.name)

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
	log.Printf("Searching %d records.", len(t.index.Items()))

	for _, meta := range t.index.Items() {
		dataObj := resultObject(t.t)
		err := t.tape.Read(meta.Point(), dataObj)

		if err != nil {
			panic(err)
		}

		if filter.Filter(dataObj) {
			if skipCount == 0 && result.Count() < pageSize {
				record := MakeRecord(meta, dataObj)
				result.add(record)
			} else {
				skipCount--
			}
		}
	}

	return result
}

func (t Table) FindFirst(filter Filterer) (Recorder, error) {
	res := t.Find(1, 1, filter)

	rator := res.GetEnumerator()
	if !rator.MoveNext() {
		return nil, errors.New("no results")
	}

	return rator.Current(), nil
}

func (t Table) Exists(filter Filterer) bool {
	_, err := t.FindFirst(filter)

	return err == nil
}

func (t Table) Create(obj Dataer) CreateSet {
	valid, err := obj.Valid()

	if err != nil {
		return CreateSet{nil, err}
	}

	if !valid {
		return CreateSet{nil, errors.New("validation failed")}
	}

	point, err := t.tape.Write(obj)

	if err != nil {
		return CreateSet{nil, err}
	}

	meta := t.index.CreateSpace(point)
	t.index.Insert(meta)

	record := MakeRecord(meta, obj)

	return CreateSet{record, nil}
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

func (t Table) Delete(key Key) error {
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

func ensureDbDirectory() {
	created := createDirectory(dbPath)

	if !created {
		log.Println("couldn't create dbPath folder")
	}
}

func ensureTableIndex(tableName, indexName string) bool {
	created := createFile(indexName)

	if !created {
		log.Println("couldn't create index for " + tableName)
	}

	return created
}

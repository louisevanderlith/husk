package husk

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

//Table controls the index and physical data tape for all records associated
type Table struct {
	t     reflect.Type
	name  string
	index Indexer
	tape  Taper
}

//NewTable returns a Table
func NewTable(obj Dataer) Tabler {
	ensureDbDirectory()

	t := reflect.TypeOf(obj).Elem()
	name := t.Name()
	idxName := getIndexName(name)
	trackName := getRecordName(name)

	ensureTableIndex(name, idxName)
	index := loadIndex(idxName)

	return Table{
		t:     t,
		name:  name,
		index: index,
		tape:  newTape(trackName),
	}
}

//FindByKey returns a Record which has the same Key
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

//Find returns a Collection of records matching the applied filter function.
func (t Table) Find(page, pageSize int, filter Filterer) Collection {
	result := NewRecordSet(page, pageSize)
	skipCount := (page - 1) * pageSize

	for _, meta := range t.index.Items() {
		dataObj := resultObject(t.t)
		err := t.tape.Read(meta.Point(), dataObj)

		if err != nil {
			panic(err)
		}

		if dataObj != nil && filter.Filter(dataObj) {
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

//FindFirst will return that first record that matches the 'filter'
func (t Table) FindFirst(filter Filterer) (Recorder, error) {
	res := t.Find(1, 1, filter)

	rator := res.GetEnumerator()
	if !rator.MoveNext() {
		return nil, errors.New("no results")
	}

	return rator.Current(), nil
}

//Exists will return true of any records match the filter.
func (t Table) Exists(filter Filterer) bool {
	_, err := t.FindFirst(filter)

	return err == nil
}

//Create adds a new data object to the collection.
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

//CreateMulti calls Create on a collection of data objects.
func (t Table) CreateMulti(objs ...Dataer) []CreateSet {
	var result []CreateSet

	for _, obj := range objs {
		set := t.Create(obj)
		result = append(result, set)
	}

	return result
}

//Update writes new data a record
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

//Delete marks the Record as Disabled and removes it from the index.
func (t Table) Delete(key Key) error {
	deleted := t.index.Delete(key)

	if !deleted {
		return errors.New("nothing deleted")
	}

	return nil
}

//Calculate does fancy stuff
func (t Table) Calculate(result interface{}, calculator Calculator) error {
	for _, meta := range t.index.Items() {
		dataObj := resultObject(t.t)
		err := t.tape.Read(meta.Point(), dataObj)

		if err != nil {
			panic(err)
		}

		if dataObj != nil {
			err = calculator.Calc(result, dataObj)

			if err != nil {
				return err
			}
		}
	}

	return nil
}

//Save writes the contents of the index
func (t Table) Save() error {
	indexName := getIndexName(t.name)

	err := write(indexName, t.index)

	if err != nil {
		return err
	}

	return nil
}

//Seed will load the seedfile into the husk database ONLY if it's empty.
func (t Table) Seed(seedfile string) error {

	if !t.Exists(Everything()) {
		result := reflect.New(reflect.SliceOf(t.t)).Interface()

		err := readJSON(seedfile, &result)

		if err != nil {
			return err
		}

		val := reflect.ValueOf(result).Elem()

		for i := 0; i < val.Len(); i++ {
			item := val.Index(i).Interface()
			t.Create(item.(Dataer))
		}
	}

	return nil
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

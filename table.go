package husk

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"

	"github.com/louisevanderlith/husk/serials"
)

//Table controls the index and physical data tape for all records associated
type Table struct {
	t     reflect.Type
	name  string
	index Indexer
	tape  Taper
}

//NewTable returns a Table
func NewTable(obj Dataer, serial Serializer) Tabler {
	ensureDbDirectory()

	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		panic("obj must not be a pointer")
	}
	gob.Register(obj)
	name := t.Name()
	idxName := getIndexName(name)
	trackName := getRecordName(name)

	ensureTableIndex(name, idxName)
	index := loadIndex(idxName)

	return Table{
		t:     t,
		name:  name,
		index: index,
		tape:  newTape(trackName, serial),
	}
}

func (t Table) Type() reflect.Type{
	return t.t
}

//FindByKey returns a Record which has the same Key
func (t Table) FindByKey(key Key) (Recorder, error) {
	var result Recorder
	meta := t.index.Get(key)

	if meta == nil {
		msg := fmt.Sprintf("key %v not found in %s", key, t.name)

		return result, errors.New(msg)
	}

	dObj := reflect.New(t.t)
	dInf := dObj.Interface()
	err := t.tape.Read(meta.Point(), &dInf)

	if err != nil {
		return nil, err
	}

	dataObj := dInf.(Dataer)
	return MakeRecord(meta, dataObj), nil
}

//Find returns a Collection of records matching the applied filter function.
func (t Table) Find(page, pageSize int, filter Filterer) (Collection, error) {
	result := NewRecordSet(page)
	skipCount := (page - 1) * pageSize

	for _, meta := range t.index.Items() {
		dObj := reflect.New(t.t)
		dInf := dObj.Interface()
		err := t.tape.Read(meta.Point(), &dInf)

		if err != nil {
			return nil, err
		}

		dataObj := dInf.(Dataer)

		if filter.Filter(dataObj) {
			if skipCount == 0 && result.Count() < pageSize {
				record := MakeRecord(meta, dataObj)
				result.add(record)
			} else {
				skipCount--
			}

			result.bean()
		}
	}

	return result, nil
}

//FindFirst will return that first record that matches the 'filter'
func (t Table) FindFirst(filter Filterer) (Recorder, error) {
	res, err := t.Find(1, 1, filter)

	if err != nil {
		return nil, err
	}

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
	if record == nil {
		return errors.New("record is empty")
	}

	data := record.Data()

	if data == nil {
		return errors.New("data is empty")
	}

	valid, err := data.Valid()

	if valid {
		meta := record.Meta()
		point, err := t.tape.Write(data)

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
		dObj := reflect.New(t.t)
		dInf := dObj.Interface()
		err := t.tape.Read(meta.Point(), &dInf)

		if err != nil {
			panic(err)
		}

		dataObj := dInf.(Dataer)

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

		jser := serials.JsonSerial{}
		byts, err := ioutil.ReadFile(seedfile)

		if err != nil {
			return nil
		}

		buffer := bytes.NewBuffer(byts)
		err = jser.Decode(buffer, &result)

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

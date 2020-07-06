package husk

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		panic("obj must not be a pointer")
	}

	gob.Register(obj)
	name := t.Name()
	idxName := getIndexName(name)
	trackName := getRecordName(name)

	idxFile, err := openFile(idxName)

	if err != nil {
		panic(err)
	}

	defer idxFile.Close()

	index, err := loadIndex(idxFile)

	if err != nil {
		panic(err)
	}

	return Table{
		t:     t,
		name:  name,
		index: index,
		tape:  newTape(trackName),
	}
}

func (t Table) Type() reflect.Type {
	return t.t
}

func (t Table) readData(m *meta) (Recorder, error) {
	dObj := reflect.New(t.t)
	dInf := dObj.Interface()
	err := t.tape.Read(m.Point(), dInf)

	if err != nil {
		return nil, err
	}

	obj := dObj.Elem().Interface().(Dataer)

	return MakeRecord(m, obj), nil
}

//FindByKey returns a Record which has the same Key
func (t Table) FindByKey(key Key) (Recorder, error) {
	var result Recorder
	meta := t.index.Get(key)

	if meta == nil {
		msg := fmt.Sprintf("key %v not found in %s", key, t.name)

		return result, errors.New(msg)
	}

	return t.readData(meta)
}

//Find returns a Collection of records matching the applied filter function.
func (t Table) Find(page, pageSize int, filter Filterer) (Collection, error) {
	result := NewRecordSet(page)
	skipCount := (page - 1) * pageSize

	for _, k := range t.index.Entries() {
		meta := t.index.Get(k)

		if meta == nil {
			continue
		}

		rec, err := t.readData(meta)

		if err != nil {
			return nil, err
		}

		if filter.Filter(rec.Data()) {
			if skipCount == 0 && result.Count() < pageSize {
				result.add(rec)
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
func (t Table) Create(obj Dataer) (Recorder, error) {
	err := obj.Valid()

	if err != nil {
		return nil, err
	}

	point, err := t.tape.Write(obj)

	if err != nil {
		return nil, err
	}

	meta := t.index.CreateSpace(point)
	t.index.Insert(meta)

	record := MakeRecord(meta, obj)

	return record, nil
}

//CreateMulti calls Create on a collection of data objects.
func (t Table) CreateMulti(objs ...Dataer) (int, error) {
	count := 0
	for _, obj := range objs {
		_, err := t.Create(obj)

		if err != nil {
			return 0, err
		}

		count++
	}

	return count, nil
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

	err := data.Valid()

	if err != nil {
		return err
	}

	point, err := t.tape.Write(data)

	if err != nil {
		return err
	}

	meta := record.Meta()
	meta.Updated(point)

	return nil
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
	for _, k := range t.index.Entries() {
		meta := t.index.Get(k)

		if meta == nil {
			continue
		}

		rec, err := t.readData(meta)

		if err != nil {
			return err
		}

		if rec.Data() == nil {
			log.Fatal("data nil for some reason")
			continue
		}

		err = calculator.Calc(result, rec.Data())

		if err != nil {
			return err
		}
	}

	return nil
}

//Save writes the contents of the index
func (t Table) Save() error {
	indexName := getIndexName(t.name)
	f, err := openFile(indexName)

	if err != nil {
		return err
	}

	defer f.Close()

	ser := gob.NewEncoder(f)
	return ser.Encode(t.index)
}

//Seed will load the seedfile into the husk database ONLY if it's empty.
func (t Table) Seed(seedfile string) error {
	if t.Exists(Everything()) {
		return nil
	}

	result := reflect.New(reflect.SliceOf(t.t)).Interface()

	byts, err := ioutil.ReadFile(seedfile)

	if err != nil {
		return err
	}

	err = json.Unmarshal(byts, result)

	if err != nil {
		return err
	}

	val := reflect.ValueOf(result).Elem()

	for i := 0; i < val.Len(); i++ {
		item := val.Index(i).Interface()
		_, err := t.Create(item.(Dataer))

		if err != nil {
			return err
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

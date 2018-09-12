package husk

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"sort"
)

type Table struct {
	t     reflect.Type
	name  string
	index *Index
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
	index := LoadIndex(path)

	return Table{
		t:     t,
		name:  name,
		index: index,
		tape:  NewTape(trackName),
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
	err := t.tape.Read(meta.Point, dataObj)

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
		err := t.tape.Read(meta.Point, dataObj)

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

func (t Table) Create(obj Dataer) CreateSet {
	set := t.CreateMulti(obj)

	if len(set) == 0 {
		return CreateSet{nil, errors.New("no records created.")}
	}

	return set[0]
}

func (t Table) CreateMulti(objs ...Dataer) []CreateSet {
	var result []CreateSet

	for _, obj := range objs {
		set := t.createRecord(obj)
		result = append(result, set)
	}

	t.index.dump(t.name)

	return result
}

func (t Table) createRecord(obj Dataer) CreateSet {
	valid, err := obj.Valid()

	if !valid {
		return CreateSet{nil, err}
	}

	nxtKey := t.index.nextKey()

	record := NewRecord(t.name, nxtKey, obj)
	meta := record.Meta()

	point, err := t.tape.Write(record.Data())

	if err == nil {
		meta.Point = point
		t.index.addMeta(meta)
	}

	return CreateSet{record, err}
}

func (t Table) Update(record Recorder) error {
	valid, err := record.Data().Valid()

	if valid {
		meta := record.Meta()
		point, err := t.tape.Write(record.Data())

		if err == nil {
			meta.Point = point
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

func (t Table) Save() {
	t.tape.Close()
	t.index.dump(t.name)
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

func (m *Index) addMeta(obj *meta) {
	key := obj.Key

	m.Values[key] = obj
	m.Keys = append(m.Keys, key)

	sort.Sort(m)
}

//Less
func (s Index) Less(i, j int) bool {
	iKey, jKey := s.Keys[i], s.Keys[j]

	return isLess(iKey, jKey)
}

func isLess(iKey, jKey Key) bool {
	stampSame := iKey.Stamp == jKey.Stamp

	if !stampSame {
		return iKey.Stamp < jKey.Stamp
	}

	return iKey.ID < jKey.ID
}

func (s Index) Swap(i, j int) {
	s.Keys[i], s.Keys[j] = s.Keys[j], s.Keys[i]
}

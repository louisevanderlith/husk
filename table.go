package husk

import (
	"errors"
	"fmt"
	"github.com/louisevanderlith/husk/collections"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/index"
	"github.com/louisevanderlith/husk/op"
	"github.com/louisevanderlith/husk/persisted"
	"github.com/louisevanderlith/husk/records"
	"github.com/louisevanderlith/husk/storage"
	"github.com/louisevanderlith/husk/storage/tape"
	"github.com/louisevanderlith/husk/validation"
	"io"
)

//Table is used to interact with records
type Table interface {
	Name() string
	//Exists confirms the existence of a record
	Exists(filter hsk.Filter) bool

	//FindByKey finds a record with a matching key.
	FindByKey(key hsk.Key) (hsk.Record, error)
	//Find looks for records that match the filter.
	Find(page, pageSize int, filter hsk.Filter) (records.Page, error)
	//FindFirst does what Find does, but will only return one record.
	FindFirst(filter hsk.Filter) (hsk.Record, error)

	//Map can modify a result set with data values
	Map(result interface{}, calculator hsk.Mapper) error

	//Create saves a new object to the database
	Create(obj validation.Dataer) (hsk.Key, error)
	//CreateMulti saves multiple records
	CreateMulti(objs ...validation.Dataer) ([]hsk.Key, error)
	//Update records changes made to a record.
	Update(key hsk.Key, obj validation.Dataer) error
	//Delete removes a record with the matching key.
	Delete(key hsk.Key) error

	//Seeds data from a json file
	//Seed(seedfile string) error
	Seed(items collections.Enumerable) error

	SaveWriter(w io.Writer) error
	Save() error
}

func NewTable(obj validation.Dataer) Table {
	store := tape.NewStore(obj, storage.GobEncoder, storage.GobDecoder)
	iFile, err := persisted.OpenIndex(store.Name())

	if err != nil {
		panic(err)
	}

	indx := index.New()
	err = persisted.LoadIndex(indx, iFile)

	if err != nil {
		panic(err)
	}

	return table{
		idx:   indx,
		store: store,
	}
}

type table struct {
	idx   hsk.Index
	store hsk.Storage
}

func (t table) SaveWriter(w io.Writer) error {
	panic("implement me")
}

func (t table) Save() error {
	return persisted.SaveIndex(t.Name(), t.idx)
}

func (t table) Name() string {
	return t.store.Name()
}

func (t table) filter(skipCount, limit int, f hsk.Filter, rec chan<- hsk.Record) {
	data := make(chan hsk.Record)

	for _, k := range t.idx.GetKeys() {
		meta := t.idx.Get(k)

		if meta == nil {
			continue
		}

		go t.store.Read(meta.Point(), data)
	}

	for i := 0; i < limit; i++ {
		obj := <-data

		if f.Filter(obj) {
			rec <- obj

			if skipCount != 0 {
				skipCount--
				continue
			}
		}
	}
}

//FindByKey returns a Record which has the same Key
func (t table) FindByKey(k hsk.Key) (hsk.Record, error) {
	meta := t.idx.Get(k)

	if meta == nil {
		return nil, fmt.Errorf("key %v not found in %s", k, t.Name())
	}

	data := make(chan hsk.Record)

	go t.store.Read(meta.Point(), data)

	return <-data, nil
}

//Find returns a Collection of records matching the applied filter function.
func (t table) Find(pageNo, pageSize int, filter hsk.Filter) (records.Page, error) {
	result := records.NewRecordPage(pageNo, pageSize)
	skipCount := (pageNo - 1) * pageSize
	totalRead := skipCount + pageSize + 1
	data := make(chan hsk.Record, totalRead)
	for _, k := range t.idx.GetKeys() {
		meta := t.idx.Get(k)

		if meta == nil {
			continue
		}

		go t.store.Read(meta.Point(), data)
	}

	for {
		rec := <-data
		if filter.Filter(rec) {
			if skipCount != 0 {
				skipCount--
				continue
			}

			if !result.Add(rec) {
				break
			}
		}
	}

	return result, nil
}

//FindFirst will return that first record that matches the 'filter'
func (t table) FindFirst(filter hsk.Filter) (hsk.Record, error) {
	for _, k := range t.idx.GetKeys() {
		meta := t.idx.Get(k)

		if meta == nil {
			continue
		}

		data := make(chan hsk.Record)
		go t.store.Read(meta.Point(), data)

		rec := <-data
		if filter.Filter(rec) {
			return rec, nil
		}
	}

	return nil, errors.New("no results")
}

//Exists will return true of any records match the filter.
func (t table) Exists(filter hsk.Filter) bool {
	_, err := t.FindFirst(filter)

	return err == nil
}

//Create adds a new data object to the collection.
func (t table) Create(obj validation.Dataer) (hsk.Key, error) {
	k, err := t.createNoSave(obj)

	if err != nil {
		return nil, err
	}

	return k, t.Save()
}

func (t table) createNoSave(obj validation.Dataer) (hsk.Key, error) {
	err := obj.Valid()

	if err != nil {
		return nil, err
	}

	point := make(chan hsk.Point)

	meta := hsk.NewMeta()
	k, err := t.idx.Add(meta)

	if err != nil {
		return nil, err
	}

	go t.store.Write(records.MakeRecord(k, obj), point)

	meta.Update(<-point)

	return k, nil
}

//CreateMulti calls Create on a collection of data objects.
func (t table) CreateMulti(objs ...validation.Dataer) ([]hsk.Key, error) {
	var result []hsk.Key

	for _, obj := range objs {
		k, err := t.createNoSave(obj)

		if err != nil {
			return nil, err
		}

		result = append(result, k)
	}

	return result, t.Save()
}

//Update writes new data a record
func (t table) Update(k hsk.Key, obj validation.Dataer) error {
	meta := t.idx.Get(k)

	if meta == nil {
		return errors.New("not item for key")
	}

	err := obj.Valid()

	if err != nil {
		return err
	}

	point := make(chan hsk.Point)

	go t.store.Write(records.MakeRecord(k, obj), point)

	meta.Update(<-point)

	return t.Save()
}

//Delete marks the Record as Disabled and removes it from the index.
func (t table) Delete(k hsk.Key) error {
	deleted := t.idx.Delete(k)

	if !deleted {
		return errors.New("nothing deleted")
	}

	return t.Save()
}

//Map allows objects to be mapped to different structures
func (t table) Map(result interface{}, calculator hsk.Mapper) error {
	for _, k := range t.idx.GetKeys() {
		meta := t.idx.Get(k)

		if meta == nil {
			continue
		}

		data := make(chan hsk.Record)
		go t.store.Read(meta.Point(), data)

		err := calculator.Map(result, <-data)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t table) Seed(items collections.Enumerable) error {
	if t.Exists(op.Everything()) {
		return nil
	}

	itor := items.GetEnumerator()
	for itor.MoveNext() {
		_, err := t.createNoSave(itor.Current().(validation.Dataer))

		if err != nil {
			return err
		}
	}

	return t.Save()
}

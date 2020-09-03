package husk

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/index"
	"github.com/louisevanderlith/husk/index/searchers"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/op"
	"github.com/louisevanderlith/husk/persisted"
	"github.com/louisevanderlith/husk/storers"
	"github.com/louisevanderlith/husk/storers/tape"
	"github.com/louisevanderlith/husk/validation"
	"io"
	"io/ioutil"
	"reflect"
)

func init() {
	gob.Register(&keys.TimeKey{})
	gob.Register(hsk.NewMeta(nil))
	gob.Register(hsk.NewPoint(0, 0))
}

func NewTable(obj validation.Dataer) hsk.Table {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		panic("obj must not be a pointer")
	}

	gob.Register(obj)

	err := persisted.CreateDirectory("db")

	if err != nil {
		panic(err)
	}

	store := tape.NewStore(t, storers.GobEncoder, storers.GobDecoder)
	return newTable(obj, store)
}

func newTable(obj validation.Dataer, store storers.Storage) hsk.Table {
	t := reflect.TypeOf(obj)
	iFile, err := persisted.OpenIndex(t.Name())

	if err != nil {
		panic(err)
	}

	indx := index.New(searchers.IndexOf)
	err = persisted.LoadIndex(indx, iFile)

	if err != nil {
		panic(err)
	}

	return table{
		objT:  t,
		idx:   indx,
		store: store,
	}
}

type table struct {
	objT  reflect.Type
	idx   hsk.Index
	store storers.Storage
}

func (t table) SaveWriter(w io.Writer) error {
	panic("implement me")
}

func (t table) Save() error {
	return persisted.SaveIndex(t.Name(), t.idx)
}

func (t table) Type() reflect.Type {
	return t.objT
}

func (t table) Name() string {
	return t.objT.Name()
}

func (t table) filter(skipCount, limit int, f hsk.Filter) (<-chan hsk.Record, error) {
	chnl := make(chan hsk.Record)
	go func() {
		for i, k := range t.idx.GetKeys() {
			meta := t.idx.Get(k)

			if meta == nil {
				continue
			}

			data := make(chan validation.Dataer)
			go t.store.Read(meta.Point(), data)

			record := hsk.MakeRecord(k, <-data)
			if f.Filter(record) {
				if skipCount != 0 {
					skipCount--
					continue
				}

				if i > limit {
					break
				}
			}
		}

		close(chnl)
	}()

	return chnl, nil
}

//FindByKey returns a Record which has the same Key
func (t table) FindByKey(k hsk.Key) (hsk.Record, error) {
	meta := t.idx.Get(k)

	if meta == nil {
		return nil, fmt.Errorf("key %v not found in %s", k, t.Name())
	}

	data := make(chan validation.Dataer)
	go t.store.Read(meta.Point(), data)

	return hsk.MakeRecord(k, <-data), nil
}

//Find returns a Collection of records matching the applied filter function.
func (t table) Find(pageNo, pageSize int, filter hsk.Filter) (hsk.Page, error) {
	result := hsk.NewRecordPage(pageNo, pageSize)
	skipCount := (pageNo - 1) * pageSize

	for _, k := range t.idx.GetKeys() {
		meta := t.idx.Get(k)

		if meta == nil {
			continue
		}

		data := make(chan validation.Dataer)
		go t.store.Read(meta.Point(), data)

		record := hsk.MakeRecord(k, <-data)

		if filter.Filter(record) {
			if skipCount != 0 {
				skipCount--
				continue
			}

			if !result.Add(record) {
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

		data := make(chan validation.Dataer)
		go t.store.Read(meta.Point(), data)

		record := hsk.MakeRecord(k, <-data)
		if filter.Filter(record) {
			return record, nil
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
	go t.store.Write(obj, point)

	meta := hsk.NewMeta(<-point)

	return t.idx.Add(meta)
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
	go t.store.Write(obj, point)

	meta.Update(<-point)

	return t.Save()
}

//Delete marks the Record as Disabled and removes it from the index.
func (t table) Delete(k hsk.Key) error {
	deleted := t.idx.Delete(k)

	if !deleted {
		return errors.New("nothing deleted")
	}

	return nil
}

//Map allows objects to be mapped to different structures
func (t table) Map(result interface{}, calculator hsk.Mapper) error {
	for _, k := range t.idx.GetKeys() {
		meta := t.idx.Get(k)

		if meta == nil {
			continue
		}

		data := make(chan validation.Dataer)
		go t.store.Read(meta.Point(), data)

		record := hsk.MakeRecord(k, <-data)

		return calculator.Map(result, record)
	}

	return nil
}

//Seed will load the seedfile into the husk database ONLY if it's empty.
func (t table) Seed(seedfile string) error {
	if t.Exists(op.Everything()) {
		return nil
	}

	result := reflect.New(reflect.SliceOf(t.Type())).Interface()

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
		item := val.Index(i).Interface().(validation.Dataer)
		_, err := t.Create(item)

		if err != nil {
			return err
		}
	}

	return nil
}

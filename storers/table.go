package storers

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"github.com/louisevanderlith/husk/op"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

func init() {
	gob.Register(&keys.TimeKey{})
	gob.Register(hsk.NewMeta(nil))
	gob.Register(hsk.NewPoint(0, 0))
}

//Table is used to interact with records
type Table interface {
	Name() string
	//Exists confirms the existence of a record
	Exists(filter op.Filter) bool

	//FindByKey finds a record with a matching key.
	FindByKey(key hsk.Key) (hsk.Record, error)
	//Find looks for records that match the filter.
	Find(page, pageSize int, filter op.Filter) (hsk.Page, error)
	//FindFirst does what Find does, but will only return one record.
	FindFirst(filter op.Filter) (hsk.Record, error)

	//Map can modify a result set with data values
	Map(result interface{}, calculator op.Mapper) error

	//Create saves a new object to the database
	Create(obj hsk.Dataer) (hsk.Key, error)
	//CreateMulti saves multiple records
	CreateMulti(objs ...hsk.Dataer) ([]hsk.Key, error)
	//Update records changes made to a record.
	Update(key hsk.Key, obj hsk.Dataer) error
	//Delete removes a record with the matching key.
	Delete(key hsk.Key) error

	//Seeds data from a json file
	Seed(seedfile string) error
	//Seeds data from a io.reader
	//SeedReader(r io.Reader) error --soon
}

type table struct {
	objT  reflect.Type
	index Index
	store Storage
}

func NewTable(obj hsk.Dataer, store Storage) Table {

	t := reflect.TypeOf(obj)
	index, err := fetchIndex(t.Name())

	if err != nil {
		panic(err)
	}

	return table{
		objT:  t,
		index: index,
		store: store,
	}
}

func (t table) Type() reflect.Type {
	return t.objT
}

func (t table) Name() string {
	return t.objT.Name()
}

func (t table) filter(skipCount, limit int, f op.Filter) (<-chan hsk.Record, error) {
	chnl := make(chan hsk.Record)
	go func() {
		for i, k := range t.index.GetKeys() {
			meta := t.index.Get(k)

			if meta == nil {
				continue
			}

			data := make(chan hsk.Dataer)
			err := t.store.Read(meta.Point(), data)

			if err != nil {
				panic(err)
			}

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
	meta := t.index.Get(k)

	if meta == nil {
		return nil, fmt.Errorf("key %v not found in %s", k, t.Name())
	}

	data := make(chan hsk.Dataer)
	err := t.store.Read(meta.Point(), data)

	if err != nil {
		return nil, err
	}

	return hsk.MakeRecord(k, <-data), nil
}

//Find returns a Collection of records matching the applied filter function.
func (t table) Find(pageNo, pageSize int, filter op.Filter) (hsk.Page, error) {
	result := hsk.NewRecordPage(pageNo, pageSize)
	skipCount := (pageNo - 1) * pageSize

	for _, k := range t.index.GetKeys() {
		meta := t.index.Get(k)

		if meta == nil {
			continue
		}

		data := make(chan hsk.Dataer)
		err := t.store.Read(meta.Point(), data)

		if err != nil {
			return nil, err
		}

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
func (t table) FindFirst(filter op.Filter) (hsk.Record, error) {
	for _, k := range t.index.GetKeys() {
		meta := t.index.Get(k)

		if meta == nil {
			continue
		}

		data := make(chan hsk.Dataer)
		err := t.store.Read(meta.Point(), data)

		if err != nil {
			return nil, err
		}

		record := hsk.MakeRecord(k, <-data)
		if filter.Filter(record) {
			return record, nil
		}
	}

	return nil, errors.New("no results")
}

//Exists will return true of any records match the filter.
func (t table) Exists(filter op.Filter) bool {
	_, err := t.FindFirst(filter)

	return err == nil
}

//Create adds a new data object to the collection.
func (t table) Create(obj hsk.Dataer) (hsk.Key, error) {
	err := obj.Valid()

	if err != nil {
		return nil, err
	}

	point, err := t.store.Write(obj)

	if err != nil {
		return nil, err
	}

	meta := hsk.NewMeta(point)
	return t.index.Add(meta)
}

//CreateMulti calls Create on a collection of data objects.
func (t table) CreateMulti(objs ...hsk.Dataer) ([]hsk.Key, error) {
	var result []hsk.Key

	for _, obj := range objs {
		k, err := t.Create(obj)

		if err != nil {
			return nil, err
		}

		result = append(result, k)
	}

	return result, nil
}

//Update writes new data a record
func (t table) Update(k hsk.Key, obj hsk.Dataer) error {
	meta := t.index.Get(k)

	if meta == nil {
		return errors.New("not item for key")
	}

	err := obj.Valid()

	if err != nil {
		return err
	}

	point, err := t.store.Write(obj)

	if err != nil {
		return err
	}

	meta.Update(point)

	return nil
	//May not need to set via index
	//return t.index.Set(k, meta)
}

//Delete marks the Record as Disabled and removes it from the index.
func (t table) Delete(k hsk.Key) error {
	deleted := t.index.Delete(k)

	if !deleted {
		return errors.New("nothing deleted")
	}

	return nil
}

//Calculate does fancy stuff
func (t table) Map(result interface{}, calculator op.Mapper) error {
	for _, k := range t.index.GetKeys() {
		meta := t.index.Get(k)

		if meta == nil {
			continue
		}

		data := make(chan hsk.Dataer)
		err := t.store.Read(meta.Point(), data)

		if err != nil {
			return err
		}

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
		item := val.Index(i).Interface().(hsk.Dataer)
		_, err := t.Create(item)

		if err != nil {
			return err
		}
	}

	return nil
}

//DestroyContents will remove the file at the given path
func DestroyContents(path string) error {
	d, err := os.Open(path)

	if err != nil {
		return err
	}

	defer d.Close()

	names, err := d.Readdirnames(-1)

	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))

		if err != nil {
			return err
		}
	}

	return nil
}

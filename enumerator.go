package husk

import (
	"errors"
)

type Enumerator interface {
	Next() (Recorder, error)
	HasNext() bool
	Any() bool
}

type recordSet []Recorder

type Enumerable struct {
	index   int
	setLen  int
	records recordSet
}

func newEnumerable(records recordSet) Enumerator {
	result := Enumerable{
		index:   0,
		setLen:  len(records),
		records: records,
	}

	return &result
}

func (e *Enumerable) add(record Recorder) {
	e.setLen++
	e.records = append(e.records, record)
}

func (e *Enumerable) Next() (result Recorder, err error) {
	if e.HasNext() {
		result = e.records[e.index]
		e.index++
	} else {
		return nil, errors.New("end of dataset")
	}

	return result, err
}

func (e *Enumerable) HasNext() bool {
	return e.index != e.setLen
}

func (e *Enumerable) Any() bool {
	return e.setLen > 0
}

func (e *Enumerable) Count() int {
	return e.setLen
}

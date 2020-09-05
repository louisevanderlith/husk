package records

import (
	"encoding/json"
	"fmt"
	"github.com/louisevanderlith/husk/collections"
	"github.com/louisevanderlith/husk/hsk"
)

//Page represents a piece of a larger dataset.
type Page interface {
	Add(record hsk.Record) bool
	collections.Enumerable
	Prev() string
	Number() int
	Next() string
	Count() int
	Any() bool
}

type page struct {
	records Collection
	number  int
	size    int
	hasMore bool
}

//NewRecordPage creates a data page for records
func NewRecordPage(number, size int) Page {
	return &page{
		number:  number,
		size:    size,
		records: NewCollection(),
		hasMore: false,
	}
}

func (s *page) GetEnumerator() collections.Iterator {
	return s.records.GetEnumerator()
}

func (s *page) Prev() string {
	prv := s.number - 1

	if prv < 1 {
		return ""
	}

	return fmt.Sprintf("%d%d", prv, s.size)
}

func (s *page) Next() string {
	if !s.hasMore {
		return ""
	}

	return fmt.Sprintf("%d%d", s.number+1, s.size)
}

func (s *page) Number() int {
	return s.number
}

//Count returns the amount of records in the collection.
func (s *page) Count() int {
	return s.records.Count()
}

//Any returns false if there are no records in the collection.
func (s *page) Any() bool {
	return s.records.Count() > 0
}

func (s *page) Add(rec hsk.Record) bool {
	if s.hasMore {
		return false
	}

	idx := s.records.Add(rec)
	s.hasMore = idx == s.size

	return idx != -1
}

//MarshalJSON returns only the 'records' instead of everything
func (s *page) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Page     int
		Previous string
		Next     string
		Length   int
		Records  Collection
	}{
		Page:     s.number,
		Previous: s.Prev(),
		Next:     s.Next(),
		Length:   s.size,
		Records:  s.records,
	})
}

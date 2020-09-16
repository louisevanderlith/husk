package records

import (
	"fmt"
	"github.com/louisevanderlith/husk/collections"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/validation"
)

//Page represents a piece of a larger dataset.
type Page interface {
	Add(record hsk.Record) bool
	collections.Enumerable
	Prev() string
	Next() string
	Count() int
	Any() bool
}

/*
	Page     int
	Previous string
	Next     string
	Length   int
	Records  Collection
*/

//NewRecordPage creates a data page for records
func NewRecordPage(t validation.Dataer, pageNo, pageSize, batchLength int) Page {
	return &page{
		Records: NewCollection(),
		Number:  pageNo,
		Size:    pageSize,
		Limit:   batchLength,
	}
}

//NewResultPage creates a page for JSONed records
func NewResultPage(t validation.Dataer) Page {
	return &page{
		Records: CollectionOf(t),
	}
}

type page struct {
	Records Collection
	Number  int
	Size    int
	Limit   int
}

func (s *page) GetEnumerator() collections.Iterator {
	return s.Records.GetEnumerator()
}

//Prev returns the Number and Size of the previous page, if available
func (s *page) Prev() string {
	prv := s.Number - 1

	if prv < 1 {
		return ""
	}

	return fmt.Sprintf("%d%d", prv, s.Size)
}

//Next returns the Number and Size of the next page, if available
func (s *page) Next() string {
	if s.Limit <= (s.Number * s.Size) {
		return ""
	}

	return fmt.Sprintf("%d%d", s.Number+1, s.Size)
}

//Count returns the amount of records in the collection.
func (s *page) Count() int {
	return s.Records.Count()
}

//Any returns false if there are no records in the collection.
func (s *page) Any() bool {
	return s.Records.Count() > 0
}

func (s *page) Add(rec hsk.Record) bool {
	if s.Count() == s.Limit {
		return false
	}

	idx := s.Records.Add(rec)
	return idx != -1
}

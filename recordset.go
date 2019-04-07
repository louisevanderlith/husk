package husk

import (
	"encoding/json"
)

//Collection is any Enumerable collection that can be counted
type Collection interface {
	Enumerable
	Count() int
	Any() bool
}

//RecordSet is a collection of Records
type RecordSet struct {
	index   int
	length  int
	records []Recorder
}

//NewRecordSet creates a collection of records.
func NewRecordSet() *RecordSet {
	return &RecordSet{
		index:  -1,
		length: 0,
	}
}

//Count returns the amount of records in the collection.
func (s *RecordSet) Count() int {
	return s.length
}

//Any returns false if there are no records in the collection.
func (s *RecordSet) Any() bool {
	return s.length > 0
}

func (s *RecordSet) add(record Recorder) {
	s.length++
	s.records = append(s.records, record)
}

//GetEnumerator returns the Enumerator used for iteration
func (s *RecordSet) GetEnumerator() Enumerator {
	s.Reset()
	return s
}

//Current returns the Record at the current position in the collection
func (s *RecordSet) Current() Recorder {
	return s.records[s.index]
}

//MoveNext is used to move to the next item in the collection
func (s *RecordSet) MoveNext() bool {
	s.index++
	return s.index < s.length
}

//Reset will set the position of "Current" the first item
func (s *RecordSet) Reset() {
	s.index = -1
	s.length = len(s.records)
}

//MarshalJSON returns only the 'records' instead of everything
func (s *RecordSet) MarshalJSON() ([]byte, error) {
	data := s.records
	return json.Marshal(data)
}

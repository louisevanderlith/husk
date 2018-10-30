package husk

import (
	"encoding/json"
	"log"
)

type Collection interface {
	Enumerable
	Count() int
	Any() bool
}

type RecordSet struct {
	index   int
	length  int
	records []Recorder
}

//NewRecordsSet creates a collection of records.
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

func (s *RecordSet) GetEnumerator() Enumerator {
	s.Reset()
	return s
}

func (s *RecordSet) Current() Recorder {
	return s.records[s.index]
}

func (s *RecordSet) MoveNext() bool {
	s.index++
	log.Printf("index: %d\n", s.index)
	return s.index < s.length
}

func (s *RecordSet) Reset() {
	s.index = -1
	s.length = len(s.records)
}

func (s *RecordSet) MarshalJSON() ([]byte, error) {
	data := s.records
	return json.Marshal(data)
}

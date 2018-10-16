package husk

import (
	"log"
)

type Collection interface {
	Enumerable
	Count() int
	Any() bool
	Add(record Recorder)
}

type RecordSet struct {
	index   int
	length  int
	records []Recorder
}

func NewRecordSet() Collection {
	return &RecordSet{
		index:  -1,
		length: 0,
	}
}

func (s *RecordSet) Count() int {
	return s.length
}

func (s *RecordSet) Any() bool {
	return s.length > 0
}

//Add adds an item to the collection. Warning! calls Reset()
func (s *RecordSet) Add(record Recorder) {
	s.length++
	s.records = append(s.records, record)
}

func (s *RecordSet) GetEnumerator() Enumerator {
	s.Reset()
	return s
}

func (s *RecordSet) Current() Recorder {
	result := s.records[s.index]

	log.Printf("RecordSET %+v\n", s)

	return result
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

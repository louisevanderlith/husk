package husk

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
	s.records = append(s.records, record)
	s.Reset()
}

func (s *RecordSet) GetEnumerator() Enumerator {
	return s
}

func (s *RecordSet) Current() Recorder {
	return s.records[s.index]
}

func (s *RecordSet) MoveNext() bool {
	s.index++

	return s.index != s.length
}

func (s *RecordSet) Reset() {
	s.index = -1
	s.length = len(s.records)
}

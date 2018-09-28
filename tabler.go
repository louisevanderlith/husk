package husk

type Tabler interface {
	FindByKey(key Key) (Recorder, error)
	Find(page, pageSize int, filter Filter) *RecordSet
	FindFirst(filter Filter) Recorder
	Exists(filter Filter) bool
	Create(objs Dataer) CreateSet
	CreateMulti(obj ...Dataer) []CreateSet
	Update(records Recorder) error
	Delete(keys Key) error
	//Writes data to disk.
	Save()
}

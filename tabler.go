package husk

type Tabler interface {
	FindByKey(key Key) (Recorder, error)
	Find(page, pageSize int, filter Filter) *RecordSet
	FindFirst(filter Filter) (Recorder, error)
	// Exists can restult in a 'true, but ...' always test for !exists
	Exists(filter Filter) (bool, error)
	Create(objs Dataer) CreateSet
	CreateMulti(obj ...Dataer) []CreateSet
	Update(records Recorder) error
	Delete(keys Key) error
	//Writes data to disk.
	Save()
}

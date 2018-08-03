package husk

type Tabler interface {
	FindByID(id int64) (Recorder, error)
	Find(page, pageSize int, filter Filter) *RecordSet
	FindFirst(filter Filter) (Recorder, error)
	// Exists can restult in a 'true, but ...' always test for !exists
	Exists(filter Filter) (bool, error)
	Create(obj Dataer) (Recorder, error)
	Update(record Recorder) error
	Delete(id int64) error
}

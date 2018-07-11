package husk

type Tabler interface {
	FindByID(id int64) (Recorder, error)
	Find(page, pageSize int, filter Filter) []Recorder
	FindFirst(filter Filter) Recorder
	Exists(filter Filter) bool
	Create(obj Dataer) (Recorder, error)
	Update(record Recorder) error
	Delete(id int64) error
}

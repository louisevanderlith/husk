package husk

type Tabler interface {
	FindByID(id int64) (Recorder, error)
	Find(page, pageSize int, filter Filter) []Recorder
	Create(obj Dataer) (Recorder, error)
	Update(record Recorder) error
	Delete(id int64) error
}

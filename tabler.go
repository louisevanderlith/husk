package husk

type Tabler interface {
	FindByKey(key *Key) (Recorder, error)
	Find(page, pageSize int, filter Filterer) Collection
	FindFirst(filter Filterer) Recorder
	Exists(filter Filterer) bool
	Create(objs Dataer) CreateSet
	CreateMulti(obj ...Dataer) []CreateSet
	Update(records Recorder) error
	Delete(keys *Key) error

	//Writes data to disk.
	Save()
}

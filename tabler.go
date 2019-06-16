package husk

//Tabler provides everything a Table should be able to do.
type Tabler interface {
	//FindByKey finds a record with a matching key.
	FindByKey(key Key) (Recorder, error)
	//Find looks for records that match the filter.
	Find(page, pageSize int, filter Filterer) Collection
	//FindFirst does what Find does, but will only return one record.
	FindFirst(filter Filterer) (Recorder, error)
	//Exists confirms the existence of a record
	Exists(filter Filterer) bool
	//Create saves a new object to the database
	Create(objs Dataer) CreateSet
	//CreateMulti saves multiple records, then commits to the database.
	CreateMulti(obj ...Dataer) []CreateSet
	//Update records changes made to a record.
	Update(records Recorder) error
	//Delete removes a record with the matching key.
	Delete(keys Key) error

	//Writes data to disk.
	Save() error

	//Seeds data from a json file
	Seed(seedfile string) error
}

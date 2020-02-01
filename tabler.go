package husk

import "reflect"

//Tabler provides everything a Table should be able to do.
type Tabler interface {
	//FindByKey finds a record with a matching key.
	FindByKey(key Key) (Recorder, error)
	//Find looks for records that match the filter.
	Find(page, pageSize int, filter Filterer) (Collection, error)
	//FindFirst does what Find does, but will only return one record.
	FindFirst(filter Filterer) (Recorder, error)
	//Calculate can modify a result set with data values
	Calculate(result interface{}, calculator Calculator) error
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
	Type() reflect.Type
}

// GetFields returns
func GetFields(t Dataer) map[string]interface{} {
	result := make(map[string]interface{})

	val := reflect.ValueOf(t)
	valType := val.Type()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := valType.Field(i)

		result[typeField.Name] = valueField.Interface()
	}

	return result
}
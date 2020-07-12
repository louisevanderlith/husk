package hsk

//Recorder is what defines a record, and what it can do
type Recorder interface {
	GetKey() Key
	//Meta() *husk.meta
	Data() Dataer
	//Set(husk.Dataer) error
}

package husk

//Recorder is what defines a record, and what it can do
type Recorder interface {
	GetKey() Key
	Meta() *meta
	Data() Dataer
	Set(Dataer) error
}

package husk

type Recorder interface {
	GetKey() *Key
	Meta() *meta
	Data() Dataer
	Set(Dataer) error
}

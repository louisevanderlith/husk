package husk

type Recorder interface {
	GetID() int64
	Meta() *meta
	Data() Dataer
	Set(Dataer) error
}

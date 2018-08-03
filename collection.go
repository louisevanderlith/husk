package husk

type Collection interface {
	Enumerable
	Count() int
	Any() bool
	Add(record Recorder)
}

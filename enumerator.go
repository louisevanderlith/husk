package husk

type Enumerator interface {
	Current() Recorder
	MoveNext() bool
	Reset()
}

package husk

type Enumerator interface {
	Current() (Recorder, error)
	MoveNext() bool
	Reset()
}

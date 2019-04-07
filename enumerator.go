package husk

//Enumerator allows iteration over collections
type Enumerator interface {
	Current() Recorder
	MoveNext() bool
	Reset()
}

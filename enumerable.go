package husk

//Enumerable is any iterable collection that has an enumerator
type Enumerable interface {
	GetEnumerator() Enumerator
}

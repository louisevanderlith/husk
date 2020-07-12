package collections

//Enumerable is any Iterable collection
type Enumerable interface {
	GetEnumerator() Iterator
}

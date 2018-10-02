package husk

//TODO: Make results enums
type Comparable interface {
	//Compare returns -1 (smaller), 0 (equal), 1 (larger)
	Compare(k Key) int8
}

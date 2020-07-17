package collections

//Comparable has the Compare Function
type Comparable interface {
	//Compare returns -1 (smaller), 0 (equal), 1 (larger)
	Compare(obj interface{}) int8
}

package collections

import "reflect"

//Iterator allows iteration over items
type Iterator interface {
	Current() interface{}
	MoveNext() bool
	Reset()
}

type itor struct {
	position int
	items    reflect.Value
}

func ReadOnlyList(slice reflect.Value) Enumerable {
	return &itor{
		position: -1,
		items:    slice,
	}
}

/*func NewIterator(items []interface{}) Iterator {
	return &itor{
		position: -1,
		items:    items,
	}
}*/

func (i *itor) GetEnumerator() Iterator {
	return i
}

func (i *itor) Current() interface{} {
	return i.items.Index(i.position).Interface()
}

func (i *itor) MoveNext() bool {
	i.position++
	return i.position < i.items.Len()
}

func (i *itor) Reset() {
	i.position = -1
}

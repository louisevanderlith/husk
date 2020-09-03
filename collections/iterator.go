package collections

//Iterator allows iteration over items
type Iterator interface {
	Current() interface{}
	MoveNext() bool
	Reset()
}

type itor struct {
	position int
	items    []interface{}
}

func NewIterator(items []interface{}) Iterator {
	return &itor{
		position: -1,
		items:    items,
	}
}

func (i *itor) Current() interface{} {
	return i.items[i.position]
}

func (i *itor) MoveNext() bool {
	i.position++
	return i.position < len(i.items)
}

func (i *itor) Reset() {
	i.position = -1
}

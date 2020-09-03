package collections

import (
	"encoding/json"
)

type List interface {
	Enumerable

	Get(index int) interface{}
	Add(rec interface{}) int
	Clear()
	Contains(rec interface{}) bool
	IndexOf(rec interface{}) int
	Insert(index int, rec interface{})
	Remove(rec interface{})
	RemoveAt(index int)
	Count() int
}

type list struct {
	items []interface{}
}

func NewList() List {
	return &list{
	}
}

func (l *list) Count() int {
	return len(l.items)
}

func (l *list) GetEnumerator() Iterator {
	return NewIterator(l.items)
}

func (l *list) Get(index int) interface{} {
	return l.items[index]
}

func (l *list) Add(rec interface{}) int {
	l.items = append(l.items, rec)

	return l.Count() - 1
}

func (l *list) Clear() {
	l.items = nil
}

func (l *list) Contains(rec interface{}) bool {
	for i := 0; i < l.Count(); i++ {
		if l.items[i] == rec {
			return true
		}
	}

	return false
}

func (l *list) IndexOf(rec interface{}) int {
	for i := 0; i < l.Count(); i++ {
		if l.items[i] == rec {
			return i
		}
	}

	return -1
}

func (l *list) Insert(index int, rec interface{}) {
	if index < 0 || index > l.Count() {
		return
	}

	l.items = append(l.items, nil)
	copy(l.items[index+1:], l.items[index:])
	l.items[index] = rec
}

func (l *list) Remove(rec interface{}) {
	l.RemoveAt(l.IndexOf(rec))
}

func (l *list) RemoveAt(index int) {
	if index >= 0 && index < l.Count() {
		for i := index; i < l.Count()-1; i++ {
			l.items[i] = l.items[i+1]
		}
	}
}

func (l *list) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.items)
}

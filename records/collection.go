package records

import (
	"github.com/louisevanderlith/husk/collections"
	"github.com/louisevanderlith/husk/hsk"
	"reflect"
)

type Collection interface {
	collections.Enumerable
	Get(index int) hsk.Record
	Add(rec hsk.Record) int
	Contains(rec hsk.Record) bool
	IndexOf(rec hsk.Record) int
	Insert(index int, rec hsk.Record)
	Remove(rec hsk.Record)
	RemoveAt(index int)
	Count() int
}

func NewCollection() Collection {
	return &collection{
	}
}

type collection struct {
	items []hsk.Record
}

func (c *collection) GetEnumerator() collections.Iterator {
	return collections.ReadOnlyList(reflect.ValueOf(c.items)).(collections.Iterator)
}

func (c *collection) Get(index int) hsk.Record {
	return c.items[index]
}

func (c *collection) Add(rec hsk.Record) int {
	i := c.Count()
	c.Insert(i, rec)
	return i
}

func (c *collection) Contains(rec hsk.Record) bool {
	return c.IndexOf(rec) != -1
}

func (c *collection) IndexOf(rec hsk.Record) int {
	for i := 0; i < c.Count(); i++ {
		if c.items[i] == rec {
			return i
		}
	}

	return -1
}

func (c *collection) Insert(index int, rec hsk.Record) {
	if index < 0 || index > c.Count() {
		return
	}

	c.items = append(c.items, nil)
	copy(c.items[index+1:], c.items[index:])
	c.items[index] = rec
}

func (c *collection) Remove(rec hsk.Record) {
	c.RemoveAt(c.IndexOf(rec))
}

func (c *collection) RemoveAt(index int) {
	if index >= 0 && index < c.Count() {
		for i := index; i < c.Count()-1; i++ {
			c.items[i] = c.items[i+1]
		}
	}
}

func (c *collection) Count() int {
	return len(c.items)
}

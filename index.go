package husk

import (
	"fmt"
	"math"
)

//Index is sorted by EPOCH Time DESC
type index struct {
	Values map[*Key]*meta
	Keys   []*Key
	Indx   int
	Total  int64
}

func loadIndex(indexName string) Indexer {
	result := &index{Values: make(map[*Key]*meta)}
	err := read(indexName, result)

	if err != nil {
		panic(err)
	}

	return result
}

// CreateSpaces generates a new Key and returns Meta
func (m *index) CreateSpace(point *Point) *meta {
	key := NewKey(m.Total)

	return NewMeta(key, point)
}

/// Create new entry in this index that maps key K to value V
func (m *index) Insert(v *meta) {
	m.Values[v.GetKey()] = v

	//key in-front
	tmp := []*Key{v.GetKey()}
	m.Keys = append(tmp, m.Keys...)

	m.Total++
}

/// Find an entry by key, returns nil of not found or not active
func (m *index) Get(k *Key) *meta {
	rec, ok := m.Values[k]

	if !ok {
		return nil
	}

	if !rec.IsActive() {
		return nil
	}

	return rec
}

/// Delete all entries of given key
func (m *index) Delete(k *Key) bool {
	idxKey := m.getKeyIndex(k)

	if idxKey == -1 {
		return false
	}

	//disable meta
	m.Values[k].Disable()

	m.Keys = append(m.Keys[:idxKey], m.Keys[idxKey+1:]...)
	fmt.Printf("disable %v :: %+v", k, m)

	m.Total--
	return true
}

func (m *index) Items() map[*Key]*meta {
	result := make(map[*Key]*meta)

	for k, meta := range m.Values {
		if meta.Active {
			result[k] = meta
		}
	}

	return result
}

func (m *index) getKeyIndex(key *Key) int {
	lft := int64(0)
	rght := m.Total - int64(1)

	for lft != rght {
		middle := int64(math.Ceil(float64((lft + rght) / 2)))
		curr := m.Keys[middle]

		//1 == greater than
		if curr.Compare(key) == 1 {
			rght = middle - 1
			continue
		}

		lft = middle
	}

	if m.Keys[lft] == key {
		return int(lft)
	}

	return -1
}

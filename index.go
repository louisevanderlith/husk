package husk

import (
	"fmt"
	"sort"
)

//Index is sorted by EPOCH Time DESC
type index struct {
	Values map[Key]*meta
	Keys   []Key
	Indx   int
	Total  int64
}

func loadIndex(indexName string) Indexer {
	result := &index{Values: make(map[Key]*meta)}
	read(indexName, result)

	return result
}

// CreateSpaces generates a new Key and returns Meta
func (m *index) CreateSpace(point *Point) *meta {
	key := NewKey(m.Total)

	return newMeta(key, point)
}

/// Create new entry in this index that maps key K to value V
func (m *index) Insert(v *meta) {
	m.Values[v.GetKey()] = v

	//key in-front
	tmp := []Key{v.GetKey()}
	m.Keys = append(tmp, m.Keys...)

	m.Total++
}

/// Find an entry by key, returns nil of not found or not active
func (m *index) Get(k Key) *meta {
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
func (m *index) Delete(k Key) bool {
	idxKey := m.getIndexOfKey(k)

	if idxKey == -1 {
		return false
	}

	meta := m.Get(k)

	if meta == nil {
		return false
	}

	meta.Disable()

	copy(m.Keys[idxKey:], m.Keys[idxKey+1:])
	m.Keys = m.Keys[:len(m.Keys)-1]

	fmt.Printf("disable %v :: %+v", k, m)

	m.Total--

	return true
}

func (m *index) Items() map[Key]*meta {
	result := make(map[Key]*meta)

	for k, meta := range m.Values {
		if meta.Active {
			result[k] = meta
		}
	}

	return result
}

func (m *index) getIndexOfKey(key Key) int {
	indx := sort.Search(len(m.Keys), func(i int) bool {
		curr := m.Keys[i]
		//Smaller or Equals, since husk is ordered by Created Date desc
		return curr.Compare(key) <= 0
	})

	if indx < len(m.Keys) && m.Keys[indx] == key {
		return indx
	}

	return -1
}

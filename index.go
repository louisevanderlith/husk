package husk

import (
	"encoding/gob"
	"os"
	"sort"
	"time"
)

//Index is sorted by EPOCH Time DESC
type index struct {
	Values map[Key]*meta
	Keys   []Key
}

func loadIndex(indexFile *os.File) (Indexer, error) {
	result := &index{Values: make(map[Key]*meta)}

	inf, err := indexFile.Stat()

	if err != nil {
		return nil, err
	}

	// new index files, won't have anything to decode
	if inf.Size() == 0 {
		return result, nil
	}

	dec := gob.NewDecoder(indexFile)
	err = dec.Decode(result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (m *index) Entries() []Key {
	return m.Keys
}

// CreateSpaces generates a new Key and returns Meta
func (m *index) CreateSpace(point *Point) *meta {
	key := m.getNextKey()

	return newMeta(key, point)
}

//getNextKey returns the next available key
func (m *index) getNextKey() Key {
	timestamp := time.Now().Unix()
	count := int64(0)
	for _, k := range m.Keys {
		if k.Stamp == timestamp {
			count++
		}
	}

	return Key{timestamp, count}
}

/// Create new entry in this index that maps key K to value V
func (m *index) Insert(v *meta) Key {
	k := v.GetKey()
	m.Values[k] = v

	//key in-front
	tmp := []Key{k}
	m.Keys = append(tmp, m.Keys...)

	return k
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

	return true
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

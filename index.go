package husk

import (
	"fmt"
)

type metaMap map[Key]*meta

//Index is sorted by EPOCH Time DESC
type Index struct {
	Values    metaMap
	Keys      []Key
	Hotfields map[string]interface{}
	indx      int
}

func newIndex() *Index {
	result := new(Index)
	result.Values = make(metaMap)
	result.Hotfields = make(map[string]interface{})

	return result
}

func (s Index) Len() int {
	return len(s.Keys)
}

func LoadIndex(indexName string) *Index {
	result := newIndex()

	err := read(indexName, result)

	if err != nil {
		panic(err)
	}

	result.Reset()

	return result
}

func (m *Index) getKeyIndex(key Key) int {
	for i, v := range m.Keys {
		if v == key {
			return i
		}
	}

	return -1
}

func (m *Index) getAt(key Key) *meta {
	meta := m.Values[key]

	if meta.Active {
		return meta
	}

	return nil
}

func (m *Index) nextKey() Key {
	nxtID := int64(1)

	if m.Len() == 0 {
		return NewKey(nxtID)
	}

	nxtID += m.Keys[0].ID

	return NewKey(nxtID)
}

func (m *Index) addMeta(obj *meta) {
	key := obj.Key

	m.Values[key] = obj
	//key in-front
	tmp := []Key{key}
	tmp = append(tmp, m.Keys...)

	m.Keys = tmp
}

func (m *Index) dump(tableName string) {
	indexName := getIndexName(tableName)

	err := write(indexName, m)

	if err != nil {
		panic(err)
	}
}

func (m *Index) disable(metaRec *meta) {
	metaRec.Disable()
	metaKey := metaRec.Key
	idxKey := m.getKeyIndex(metaKey)
	m.Keys = append(m.Keys[:idxKey], m.Keys[idxKey+1:]...)

	fmt.Printf("disable %v :: %+v", metaRec.Active, m)
}

func (m *Index) Current() *meta {
	k := m.Keys[m.indx]
	curr := m.Values[k]

	return curr
}

func (m *Index) MoveNext() bool {
	m.indx--

	if m.indx == 0 {
		m.Reset()
		return false
	}

	return true
}

func (m *Index) Reset() {
	m.indx = m.Len()
}

package husk

import (
	"fmt"
	"sort"
	"time"
)

//Index is sorted by EPOCH Time DESC
type Index struct {
	Values    map[Key]*meta
	Keys      []Key
	Hotfields map[string]interface{}
	indx      int
}

func newIndex() *Index {
	result := new(Index)
	result.Values = make(map[Key]*meta)
	result.Hotfields = make(map[string]interface{})

	return result
}

func (s Index) Len() int {
	return len(s.Keys)
}

//Less
func (s Index) Less(i, j int) bool {
	iKey, jKey := s.Keys[i], s.Keys[j]

	return isLess(iKey, jKey)
}

func isLess(iKey, jKey Key) bool {
	stampSame := iKey.Stamp == jKey.Stamp

	if !stampSame {
		return iKey.Stamp < jKey.Stamp
	}

	return iKey.ID < jKey.ID
}

func (s Index) Swap(i, j int) {
	s.Keys[i], s.Keys[j] = s.Keys[j], s.Keys[i]
}

func LoadIndex(indexName string) *Index {
	result := newIndex()

	err := read(indexName, result)

	if err != nil {
		panic(err)
	}

	result.indx = -1

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

	if meta != nil && meta.Active {
		return meta
	}

	return nil
}

func (m *Index) nextKey() Key {
	stamp := time.Now().UnixNano()
	nxtID := int64(1)

	if len(m.Keys) == 0 {
		return Key{stamp, nxtID}
	}

	top := m.Keys[0]
	nxtID += top.ID

	return Key{stamp, nxtID}
}

func (m *Index) addMeta(obj *meta) {
	key := obj.Key

	m.Values[key] = obj
	m.Keys = append(m.Keys, key)

	sort.Sort(m)
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

	kIdx := m.getKeyIndex(metaRec.Key)
	m.Keys = append(m.Keys[:kIdx], m.Keys[kIdx+1:]...)

	fmt.Printf("disable %b :: %+v", metaRec.Active, m)
}

func (m *Index) Current() *meta {
	k := m.Keys[m.indx]

	return m.Values[k]
}

func (m *Index) MoveNext() bool {
	m.indx++

	if m.indx == len(m.Keys) {
		m.Reset()
		return false
	}

	return true
}

func (m *Index) Reset() {
	m.indx = -1
}

package storers

import (
	"github.com/louisevanderlith/husk/hsk"
	"sort"
	"sync"
)

//Indexer is used to manage data and where it's located
type Indexer interface {
	/// Create new entry in this index that maps value.Key K to value V
	Add(v *hsk.Point) hsk.Meta
	Set(k hsk.Key, v hsk.Meta)
	/// Find an entry by key
	Get(k hsk.Key) hsk.Meta
	/// Delete all entries of given key
	Delete(k hsk.Key) bool

	GetKeys() []hsk.Key
}

func NewIndex() Indexer {
	return &index{
		rwm:    sync.RWMutex{},
		Values: make(map[hsk.Key]hsk.Meta),
		Keys:   nil,
	}
}

type index struct {
	rwm    sync.RWMutex
	Values map[hsk.Key]hsk.Meta
	Keys   []hsk.Key
}

func (i *index) Add(p *hsk.Point) hsk.Meta {
	k := i.getNextKey()

	m := hsk.NewMeta(k, p)
	i.Set(k, m)

	return m
}

func (i *index) Set(k hsk.Key, v hsk.Meta) {
	i.rwm.Lock()
	i.Values[k] = v
	i.rwm.Unlock()

	//key in-front
	tmp := []hsk.Key{k}
	i.Keys = append(tmp, i.Keys...)
}

func (i *index) Get(k hsk.Key) hsk.Meta {
	i.rwm.RLock()
	rec, ok := i.Values[k]
	defer i.rwm.RUnlock()
	if !ok {
		return nil
	}

	if !rec.IsActive() {
		return nil
	}

	return rec
}

func (i *index) Delete(k hsk.Key) bool {
	idxKey := i.getIndexOfKey(k)

	if idxKey == -1 {
		return false
	}

	meta := i.Get(k)

	if meta == nil {
		return false
	}

	meta.Disable()

	copy(i.Keys[idxKey:], i.Keys[idxKey+1:])
	i.Keys = i.Keys[:len(i.Keys)-1]

	i.rwm.Lock()
	defer i.rwm.Unlock()
	delete(i.Values, k)

	return true
}

func (i *index) GetKeys() []hsk.Key {
	return i.Keys
}

func (i *index) getIndexOfKey(key hsk.Key) int {
	indx := sort.Search(len(i.Keys), func(n int) bool {
		curr := i.Keys[n]
		//Smaller or Equals, since husk is ordered by Created Date desc
		return curr.Compare(key) <= 0
	})

	if indx < len(i.Keys) && i.Keys[indx] == key {
		return indx
	}

	return -1
}

//getNextKey returns the next available key
func (i *index) getNextKey() hsk.Key {
	k := hsk.NewKey(0)

	i.rwm.RLock()
	defer i.rwm.RUnlock()

	for {
		_, hasKey := i.Values[k]

		if !hasKey {
			return k
		}

		k.ID++
	}

	return hsk.CrazyKey()
}

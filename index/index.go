package index

import (
	"encoding/gob"
	"errors"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"sort"
	"sync"
	"time"
)

func init() {
	gob.Register(hsk.NewMeta())
	gob.Register(hsk.NewPoint(0, 0))
}

func New() hsk.Index {
	return &index{
		rwm:    sync.RWMutex{},
		Values: make(map[hsk.Key]hsk.Meta),
		Keys:   nil,
	}
}

type searchFunc func(n int, f func(int) bool) int

type index struct {
	search searchFunc
	rwm    sync.RWMutex
	Values map[hsk.Key]hsk.Meta
	Keys   []hsk.Key
}

func (i *index) Add(m hsk.Meta) (hsk.Key, error) {
	k := i.getNextKey()

	if k == nil {
		return nil, errors.New("unable to generate key")
	}

	idx := search(i.Keys, k)

	i.Keys = append(i.Keys, nil)
	copy(i.Keys[idx+1:], i.Keys[idx:])
	i.Keys[idx] = k

	i.rwm.Lock()
	defer i.rwm.Unlock()

	i.Values[k] = m

	return k, nil
}

func (i *index) Set(k hsk.Key, v hsk.Meta) error {
	idx := i.IndexOf(k)

	if idx == -1 {
		return errors.New("unable to find key")
	}

	i.rwm.Lock()
	defer i.rwm.Unlock()

	i.Values[k] = v

	return nil
}

func (i *index) Get(k hsk.Key) hsk.Meta {
	idx := i.IndexOf(k)

	if idx == -1 {
		return nil
	}

	i.rwm.RLock()
	defer i.rwm.RUnlock()

	rec, ok := i.Values[k]

	if !ok {
		return nil
	}

	if !rec.IsActive() {
		return nil
	}

	return rec
}

func (i *index) Delete(k hsk.Key) bool {
	idx := i.IndexOf(k)

	if idx == -1 {
		return false
	}

	copy(i.Keys[idx:], i.Keys[idx+1:])
	i.Keys = i.Keys[:len(i.Keys)-1]

	i.rwm.Lock()
	defer i.rwm.Unlock()

	meta := i.Values[k]

	if meta == nil {
		return false
	}

	meta.Disable()

	return true
}

func (i *index) GetKeys() []hsk.Key {
	return i.Keys
}

func (i *index) IndexOf(k hsk.Key) int {
	idx := search(i.Keys, k)

	if idx < len(i.Keys) && i.Keys[idx].Compare(k) == 0 {
		return idx
	}

	return -1
}

func search(keys []hsk.Key, k hsk.Key) int {
	return sort.Search(len(keys), func(n int) bool {
		curr := keys[n]
		//Smaller or Equals, since husk is ordered by Created Date desc
		comp := curr.Compare(k)
		return comp <= 0
	})
}

//getNextKey returns the next available key
func (i *index) getNextKey() hsk.Key {
	timestamp := time.Now().Unix()
	id := int64(0)

	k := keys.NewKeyWithTime(timestamp, id)
	for {
		if i.IndexOf(k) == -1 {
			return k
		}

		id++
		k = keys.NewKeyWithTime(timestamp, id)
	}

	return nil
}

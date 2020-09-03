package index

import (
	"errors"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"sync"
	"time"
)

func New(search searchFunc) hsk.Index {
	return &index{
		search: search,
		rwm:    sync.RWMutex{},
		Values: make(map[int]hsk.Meta),
		Keys:   nil,
	}
}

type searchFunc func(keys []hsk.Key, key hsk.Key) int

type index struct {
	search searchFunc
	rwm    sync.RWMutex
	Values map[int]hsk.Meta
	Keys   []hsk.Key
}

func (i *index) Add(m hsk.Meta) (hsk.Key, error) {
	k := i.getNextKey()

	if k == nil {
		return nil, errors.New("unable to generate key")
	}

	i.rwm.Lock()
	defer i.rwm.Unlock()

	i.Values[len(i.Keys)] = m

	//key in-front
	tmp := []hsk.Key{k}
	i.Keys = append(tmp, i.Keys...)

	return k, nil
}

func (i *index) Set(k hsk.Key, v hsk.Meta) error {
	idx := i.IndexOf(k)

	if idx == -1 {
		return errors.New("unable to find key")
	}

	i.rwm.Lock()
	defer i.rwm.Unlock()

	i.Values[idx] = v

	return nil
}

func (i *index) Get(k hsk.Key) hsk.Meta {
	idx := i.IndexOf(k)

	if idx == -1 {
		return nil
	}

	i.rwm.RLock()
	defer i.rwm.RUnlock()

	rec, ok := i.Values[idx]

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

	meta := i.Values[idx]

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
	return i.search(i.Keys, k)
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

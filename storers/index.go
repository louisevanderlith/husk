package storers

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/keys"
	"os"
	"sort"
	"sync"
	"time"
)

//Index is used to manage data and where it's located
type Index interface {
	/// Create new entry in this index that maps value.Key K to value V
	Add(m hsk.Meta) (hsk.Key, error)
	Set(k hsk.Key, v hsk.Meta) error
	/// Find an entry by key
	Get(k hsk.Key) hsk.Meta
	IndexOf(k hsk.Key) int
	/// Delete all entries of given key
	Delete(k hsk.Key) bool

	GetKeys() []hsk.Key
}

func NewIndex() Index {
	return &index{
		rwm:    sync.RWMutex{},
		Values: make(map[int]hsk.Meta),
		Keys:   nil,
	}
}

type index struct {
	rwm    sync.RWMutex
	Values map[int]hsk.Meta
	Keys   []hsk.Key
}

func (i *index) Add(m hsk.Meta) (hsk.Key, error) {
	k := i.getNextKey()

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
	delete(i.Values, idx)

	meta := i.Get(k)

	if meta == nil {
		return false
	}

	meta.Disable()

	return true
}

func (i *index) GetKeys() []hsk.Key {
	return i.Keys
}

func (i *index) IndexOf(key hsk.Key) int {
	idx := sort.Search(len(i.Keys), func(n int) bool {
		curr := i.Keys[n]
		//Smaller or Equals, since husk is ordered by Created Date desc
		comp := curr.Compare(key)
		return comp <= 0
	})

	if idx < len(i.Keys) && i.Keys[idx].Compare(key) == 0 {
		return idx
	}

	return -1
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

	return keys.CrazyKey()
}

func fetchIndex(name string) (Index, error) {
	idxFile, err := openIndex(name)

	if err != nil {
		return nil, err
	}

	defer idxFile.Close()

	return loadIndex(idxFile)
}

func openIndex(name string) (*os.File, error) {
	err := createDirectory("db")

	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("db/%s.Index.husk", name)
	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
}

func loadIndex(indexFile *os.File) (Index, error) {
	inf, err := indexFile.Stat()

	if err != nil {
		return nil, err
	}

	result := NewIndex()

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

func saveIndex(name string, i Index) error {
	idxFile, err := openIndex(name)

	if err != nil {
		return err
	}

	defer idxFile.Close()

	serial := gob.NewEncoder(idxFile)
	return serial.Encode(i)
}

func createDirectory(folderPath string) error {
	_, err := os.Stat(folderPath)

	if err == nil {
		return nil
	}

	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		return os.Mkdir(folderPath, 0755)
	}

	return nil
}

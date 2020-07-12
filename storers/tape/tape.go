package tape

import (
	"encoding/gob"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
	"log"
	"os"
	"reflect"
)

func NewTable(obj hsk.Dataer) storers.Table {
	created := createDirectory(dbPath)

	if !created {
		log.Println("couldn't create dbPath folder")
	}

	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		panic("obj must not be a pointer")
	}

	gob.Register(obj)

	index, err := fetchIndex(t.Name())

	if err != nil {
		panic(err)
	}

	return storers.NewTable(t, index, newStore(t))
}

func fetchIndex(name string) (storers.Indexer, error) {
	idxName := getIndexPath(name)
	idxFile, err := openFile(idxName)

	if err != nil {
		return nil, err
	}

	defer idxFile.Close()

	return loadIndex(idxFile)
}

func loadIndex(indexFile *os.File) (storers.Indexer, error) {
	inf, err := indexFile.Stat()

	if err != nil {
		return nil, err
	}

	result := storers.NewIndex()

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

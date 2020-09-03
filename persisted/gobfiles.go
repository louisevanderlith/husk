package persisted

import (
	"encoding/gob"
	"fmt"
	"github.com/louisevanderlith/husk/hsk"
	"os"
)

func LoadIndex(result hsk.Index, indexFile *os.File) error {
	inf, err := indexFile.Stat()

	if err != nil {
		return err
	}

	// new index files, won't have anything to decode
	if inf.Size() == 0 {
		return nil
	}

	dec := gob.NewDecoder(indexFile)
	err = dec.Decode(result)

	if err != nil {
		return err
	}

	return nil
}

func OpenIndex(name string) (*os.File, error) {
	err := CreateDirectory("db")

	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("db/%s.Index.husk", name)
	return os.Create(filePath)
}

func SaveIndex(name string, i hsk.Index) error {
	idxFile, err := OpenIndex(name)

	if err != nil {
		return err
	}

	defer idxFile.Close()

	serial := gob.NewEncoder(idxFile)
	return serial.Encode(i)
}

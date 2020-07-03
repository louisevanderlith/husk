package husk

import (
	"encoding/gob"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	registerGobTypes()
}

func registerGobTypes() {
	gob.Register(index{})
	gob.Register(meta{})
}

func write(filePath string, data interface{}) error {
	f, err := openFile(filePath)

	if err != nil {
		return err
	}

	ser := gob.NewEncoder(f)
	return ser.Encode(data)
}

func read(f *os.File, result interface{}) error {
	dec := gob.NewDecoder(f)
	err := dec.Decode(result)

	if err == io.EOF {
		return nil
	}

	return err
}

func getDBIndexFiles() map[string]string {
	result := make(map[string]string)

	files, err := ioutil.ReadDir(dbPath)

	if err != nil {
		return result
	}

	for _, v := range files {
		path := v.Name()

		if strings.HasSuffix(path, ".index.husk") {
			key := cleanTableName(path)
			result[key] = path
		}
	}

	return result
}

func openFile(filePath string) (*os.File, error) {
	_, err := os.Stat(filePath)

	if os.IsNotExist(err) {
		return os.Create(filePath)
	} else if err != nil {
		return nil, err
	}

	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)
}

func createDirectory(folderPath string) bool {
	_, err := os.Stat(folderPath)
	notExist := os.IsNotExist(err)

	if notExist {
		err = os.Mkdir(folderPath, 0755)

		notExist = err != nil
	}

	return !notExist
}

//DestroyContents will remove the file at the given path
func DestroyContents(path string) error {
	d, err := os.Open(path)

	if err != nil {
		return err
	}

	defer d.Close()

	names, err := d.Readdirnames(-1)

	if err != nil {
		return err
	}

	for _, name := range names {
		err = os.RemoveAll(filepath.Join(path, name))

		if err != nil {
			return err
		}
	}

	return nil
}

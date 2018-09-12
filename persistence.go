package husk

import (
	"bytes"
	"encoding/gob"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func registerGobTypes() {
	gob.Register(Index{})
	gob.Register(dataMap{})
	gob.Register(meta{})
	gob.Register(Record{})
}

func write(filePath string, data interface{}) error {
	created := createFile(filePath)

	if !created {
		return errors.New("unable to create " + filePath)
	}

	bytes, err := toBytes(data)

	if err == nil {
		err = ioutil.WriteFile(filePath, bytes, 0644)
	}

	return err
}

func read(filePath string, result interface{}) error {
	byts, err := ioutil.ReadFile(filePath)

	if err != nil {
		return err
	}

	if len(byts) != 0 {
		buffer := bytes.NewBuffer(byts)
		err = gob.NewDecoder(buffer).Decode(result)
	}

	return err
}

func getDBIndexFiles() map[string]string {
	result := make(map[string]string)

	files, err := ioutil.ReadDir(dbPath)

	if err != nil {
		log.Print("getDBIndexFiles ", err)
		return result
	}

	for _, v := range files {
		path := v.Name()

		if strings.Contains(path, "index.oli") {
			key := cleanTableName(path)
			result[key] = path
		}
	}

	return result
}

func toBytes(obj interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	err := gob.NewEncoder(buffer).Encode(obj)

	return buffer.Bytes(), err
}

func createFile(filePath string) bool {
	_, err := os.Stat(filePath)
	notexist := os.IsNotExist(err)

	if notexist {
		var file, err = os.Create(filePath)

		if err != nil {
			log.Println(err)
			return false
		}

		notexist = false
		defer file.Close()
	}

	return !notexist
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

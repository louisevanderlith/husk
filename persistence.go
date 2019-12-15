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

	"github.com/louisevanderlith/husk/serials"
)

func init() {
	registerGobTypes()
}

func registerGobTypes() {
	gob.Register(index{})
	gob.Register(meta{})
}

func write(filePath string, data interface{}) error {
	created := createFile(filePath)

	if !created {
		return errors.New("unable to create " + filePath)
	}

	ser := serials.GobSerial{}
	bytes, err := ser.Encode(data) //toBytes(data)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, bytes, 0644)
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

		if strings.HasSuffix(path, ".index.husk") {
			key := cleanTableName(path)
			result[key] = path
		}
	}

	return result
}

func createFile(filePath string) bool {
	_, err := os.Stat(filePath)
	notexist := os.IsNotExist(err)

	if notexist {
		var file, err = os.Create(filePath)

		if err != nil {
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

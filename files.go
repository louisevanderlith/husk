package husk

import (
	"fmt"
	"log"
	"strings"
)

const (
	dbPath       string = "db"
	dataPath     string = "db/%s.data.husk"
	indexPath    string = "db/%s.index.husk"
	indexPattern string = "db/*.index.husk"
)

func ensureDbDirectory() {
	created := createDirectory(dbPath)

	if !created {
		log.Println("couldn't create dbPath folder")
	}
}

func getRecordName(tableName string) string {
	return fmt.Sprintf(dataPath, tableName)
}

func getIndexName(tableName string) string {
	return fmt.Sprintf(indexPath, tableName)
}

func cleanTableName(indexName string) string {
	var result string

	result = strings.Replace(indexName, ".index.husk", "", 1)

	return result
}

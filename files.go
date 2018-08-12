package husk

import (
	"fmt"
	"log"
	"strings"
)

const (
	dbPath       string = "db"
	recordPath   string = "db/%s.%v.husk"
	indexPath    string = "db/%s.index.husk"
	indexPattern string = "db/*.index.husk"
)

func ensureDbDirectory() {
	created := createDirectory(dbPath)

	if !created {
		log.Println("couldn't create dbPath folder")
	}
}

func getRecordName(tableName string, id int64) string {
	return fmt.Sprintf(recordPath, tableName, id)
}

func getIndexName(tableName string) string {
	return fmt.Sprintf(indexPath, tableName)
}

func cleanTableName(indexName string) string {
	var result string

	result = strings.Replace(indexName, ".index.husk", "", 1)

	return result
}

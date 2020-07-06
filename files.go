package husk

import (
	"fmt"
	"strings"
)

const (
	dbPath    string = "db"
	dataPath  string = "db/%s.Data.husk"
	indexPath string = "db/%s.index.husk"
)

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

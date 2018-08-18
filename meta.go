package husk

import (
	"time"
)

type meta struct {
	Key         Key
	Active      bool
	FileName    string
	LastUpdated time.Time
}

func NewMeta(tableName string, key Key) (result *meta) {
	fileName := getRecordName(tableName, key)
	created := createFile(fileName)

	if created {
		result = &meta{
			Key:         key,
			Active:      true,
			FileName:    fileName,
			LastUpdated: time.Now(),
		}
	}

	return result
}

func (m *meta) Disable() {
	m.Active = false

}

func (m *meta) Updated() {
	m.LastUpdated = time.Now()
}

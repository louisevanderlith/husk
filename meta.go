package husk

import (
	"time"
)

type meta struct {
	ID          int64
	Active      bool
	CreateDate  time.Time
	FileName    string
	LastUpdated time.Time
}

func NewMeta(tableName string, id int64) (result *meta) {
	fileName := getRecordName(tableName, id)
	created := createFile(fileName)

	if created {
		result = &meta{
			ID:          id,
			Active:      true,
			CreateDate:  time.Now(),
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

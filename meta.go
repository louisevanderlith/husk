package husk

import (
	"time"
)

type meta struct {
	Key         Key
	Active      bool
	Point       *Point
	LastUpdated time.Time
}

func NewMeta(key Key) (result *meta) {
	return &meta{
		Key:         key,
		Active:      true,
		LastUpdated: time.Now(),
	}
}

func (m *meta) Disable() {
	m.Active = false

}

func (m *meta) Updated() {
	m.LastUpdated = time.Now()
}

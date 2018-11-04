package husk

import (
	"time"
)

type meta struct {
	Key     Key
	Active  bool
	Pointer *Point
	Changed time.Time
}

func NewMeta(key Key, point *Point) *meta {
	return &meta{
		Key:     key,
		Active:  true,
		Pointer: point,
		Changed: time.Now(),
	}
}

func (m *meta) Disable() {
	m.Active = false
}

func (m *meta) Updated(p *Point) {
	m.Pointer = p
	m.Changed = time.Now()
}

func (m *meta) GetKey() Key {
	return m.Key
}

func (m *meta) IsActive() bool {
	return m.Active
}

func (m *meta) Point() *Point {
	return m.Pointer
}

func (m *meta) LastUpdated() time.Time {
	return m.Changed
}

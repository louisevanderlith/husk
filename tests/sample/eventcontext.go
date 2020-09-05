package sample

import (
	"github.com/louisevanderlith/husk"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/op"
	"github.com/louisevanderlith/husk/records"
)

type EventContext interface {
	DeleteEvent(k hsk.Key) error
	UpdateEvent(k hsk.Key, obj Event) error
	GetEvent(k hsk.Key) (hsk.Record, error)
	FindEvents(page, size int) (records.Page, error)
	FindEventsByType(page, size int, name string) (records.Page, error)
	HasEvents() bool
	CreateEvent(t string, rel hsk.Key) (hsk.Key, error)
}

//Returns a new Journal Database
func NewEventContext() EventContext {
	//return husk.NewDatabase()
	return eventcontext{
		Events: husk.NewTable(Event{}),
	}
}

type eventcontext struct {
	Events husk.Table
}

func (e eventcontext) DeleteEvent(k hsk.Key) error {
	return e.Events.Delete(k)
}

func (e eventcontext) UpdateEvent(k hsk.Key, obj Event) error {
	return e.Events.Update(k, obj)
}

func (e eventcontext) GetEvent(k hsk.Key) (hsk.Record, error) {
	return e.Events.FindByKey(k)
}

func (e eventcontext) FindEvents(page, size int) (records.Page, error) {
	return e.Events.Find(page, size, op.Everything())
}

func (e eventcontext) FindEventsByType(page, size int, name string) (records.Page, error) {
	return e.Events.Find(page, size, ByType(name))
}

func (e eventcontext) HasEvents() bool {
	return e.Events.Exists(op.Everything())
}

func (e eventcontext) CreateEvent(t string, rel hsk.Key) (hsk.Key, error) {
	return e.Events.Create(Event{
		Type:     t,
		Relation: rel,
	})
}

func (e eventcontext) Save() error {
	return nil
}

func (e eventcontext) Seed() error {
	return nil
}

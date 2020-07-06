package sample

import (
	"github.com/louisevanderlith/husk"
)

type journalFilter func(obj Journal) bool

func (f journalFilter) Filter(obj husk.Dataer) bool {
	return f(obj.(Journal))
}

func ByPublisher(name string) journalFilter {
	return func(obj Journal) bool {
		return obj.Entry.Publisher == name
	}
}

func ByObject(parm Journal) journalFilter {
	return func(obj Journal) bool {
		//More fields can be added
		return len(parm.Entry.Title) == 0 || obj.Entry.Title == parm.Entry.Title
	}
}

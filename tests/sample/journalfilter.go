package sample

import (
	"github.com/louisevanderlith/husk/hsk"
)

type journalFilter func(obj Journal) bool

func (f journalFilter) Filter(obj hsk.Dataer) bool {
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

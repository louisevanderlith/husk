package sample

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/op"
)

func ByPublisher(name string) op.FilterFunc {
	return func(rec hsk.Record) bool {
		obj := rec.Data().(Journal)
		return obj.Entry.Publisher == name
	}
}

func ByObject(parm Journal) op.FilterFunc {
	return func(rec hsk.Record) bool {
		obj := rec.Data().(Journal)
		//More fields can be added

		return len(parm.Entry.Title) == 0 || obj.Entry.Title == parm.Entry.Title
	}
}

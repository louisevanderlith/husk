package sample

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/op"
)

func ByType(t string) op.FilterFunc {
	return func(rec hsk.Record) bool {
		obj := rec.GetValue().(Event)
		return obj.Type == t
	}
}

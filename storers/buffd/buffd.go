package buffd

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
	"reflect"
)

func NewTable(obj hsk.Dataer) storers.Table {
	t := reflect.TypeOf(obj)
	return storers.NewTable(t, storers.NewIndex(), newStore())
}

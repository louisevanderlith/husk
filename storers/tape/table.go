package tape

import (
	"encoding/gob"
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
	"reflect"
)

func init() {
	gob.Register(tapeStore{})
}

func NewTable(obj hsk.Dataer) storers.Table {
	t := reflect.TypeOf(obj)

	if t.Kind() == reflect.Ptr {
		panic("obj must not be a pointer")
	}

	gob.Register(obj)

	return storers.NewTable(obj, newStore(t))
}

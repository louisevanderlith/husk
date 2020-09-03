package op

import (
	"github.com/louisevanderlith/husk/hsk"
)

type MapperFunc func(result interface{}, obj hsk.Record) error

func (f MapperFunc) Map(result interface{}, obj hsk.Record) error {
	return f(result, obj)
}

func RowCount() MapperFunc {
	return func(result interface{}, obj hsk.Record) error {
		count := result.(*int64)

		*count++

		return nil
	}
}

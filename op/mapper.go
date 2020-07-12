package op

import (
	"github.com/louisevanderlith/husk/hsk"
)

//Calculator updates the result set with values from data
type Mapper interface {
	Map(result interface{}, obj hsk.Dataer) error
}

type mapperFunc func(result interface{}, obj hsk.Dataer) error

func (f mapperFunc) Map(result interface{}, obj hsk.Dataer) error {
	return f(result, obj)
}

func RowCount() mapperFunc {
	return func(result interface{}, obj hsk.Dataer) error {
		count := result.(*int64)

		*count++

		return nil
	}
}

package husk

//Calculator updates the result set with values from data
type Calculator interface {
	Calc(result interface{}, obj Dataer) error
}

type calcer func(result interface{}, obj Dataer) error

func (c calcer) Calc(result interface{}, obj Dataer) error {
	return c(result, obj)
}

func RowCount() calcer {
	return func(result interface{}, obj Dataer) error {
		count := result.(*int64)

		*count++

		return nil
	}
}

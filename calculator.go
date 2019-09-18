package husk

//Calculator updates the result set with values from data
type Calculator interface {
	Calc(result interface{}, obj Dataer) error
}

type calcer func(obj Dataer) bool

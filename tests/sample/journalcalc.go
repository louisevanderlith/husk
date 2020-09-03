package sample

import (
	"github.com/louisevanderlith/husk/validation"
)

type journalCalc func(result interface{}, obj Journal) error

func (f journalCalc) Calc(result interface{}, obj validation.Dataer) error {
	return f(result, obj.(Journal))
}

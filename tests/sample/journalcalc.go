package sample

import (
	"github.com/louisevanderlith/husk"
)

type journalCalc func(result interface{}, obj Journal) error

func (f journalCalc) Calc(result interface{}, obj husk.Dataer) error {
	return f(result, obj.(Journal))
}

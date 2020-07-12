package sample

import (
	"github.com/louisevanderlith/husk/hsk"
)

type journalCalc func(result interface{}, obj Journal) error

func (f journalCalc) Calc(result interface{}, obj hsk.Dataer) error {
	return f(result, obj.(Journal))
}

package sample

import (
	"github.com/louisevanderlith/husk"
)

type personCalc func(result interface{}, obj *Person) error

func (f personCalc) Calc(result interface{}, obj husk.Dataer) error {
	return f(result, obj.(*Person))
}

func SumBalance() personCalc {
	return func(result interface{}, obj *Person) error {
		answ := float32(0)
		for _, acc := range obj.Accounts {
			answ += acc.Balance
		}

		totl := result.(*float32)
		*totl += answ

		return nil
	}
}

func LowestBalance() personCalc {
	min := float32(9999999)
	return func(result interface{}, obj *Person) error {
		answ := float32(0)
		for _, acc := range obj.Accounts {
			answ += acc.Balance
		}

		if answ < min {
			min = answ
			n := result.(*string)
			*n = obj.Name
		}

		return nil
	}
}

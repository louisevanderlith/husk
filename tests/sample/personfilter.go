package sample

import (
	"log"

	"github.com/louisevanderlith/husk"
)

type personFilter func(obj Person) bool

func (f personFilter) Filter(obj husk.Dataer) bool {
	return f(obj.(Person))
}

func ByName(name string) personFilter {
	return func(obj Person) bool {
		log.Printf("%+v\n", obj)
		return obj.Name == name
	}
}

func SameBalance(balance float32) personFilter {
	return func(obj Person) bool {
		for _, v := range obj.Accounts {
			if v.Balance == balance {
				return true
			}
		}

		return false
	}
}

func ByNameAndAge(name string, age int) personFilter {
	return func(obj Person) bool {
		return obj.Name == name && obj.Age == age
	}
}

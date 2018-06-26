package sample

import "github.com/louisevanderlith/husk"

type Person struct {
	Name string
	Age  int
}

func GetByName(p Person, name string) husk.Filter {
	return func(obj husk.Dataer) bool {
		return p.Name == name
	}
}

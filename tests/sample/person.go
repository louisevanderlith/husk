package sample

import "github.com/louisevanderlith/husk"

type Person struct {
	Name string `hsk:"size(50)"`
	Age  int
}

func (o Person) Valid() (bool, error) {
	return husk.ValidateStruct(o)
}

func GetByName(p Person, name string) husk.Filter {
	return func(obj husk.Dataer) bool {
		return p.Name == name
	}
}

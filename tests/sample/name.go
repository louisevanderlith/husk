package sample

type Name struct {
	string
}

func (o Name) Valid() (bool, error) {
	return len(o.string) != 0, nil
}

package husk

type Dataer interface {
	Valid() (bool, error)
}

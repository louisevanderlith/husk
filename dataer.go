package husk

//Dataer is the primary interface that any "model" should implement
//"Models" are data objects used to store and structure records in tabes.
type Dataer interface {
	Valid() (bool, error)
}

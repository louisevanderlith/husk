package hsk

// Filter used to filter records while searching.
type Filter interface {
	// Enables a single central function to cast from husk.Dataer to <Entity>
	Filter(obj Record) bool
}

//Mapper updates the result set with values from data
type Mapper interface {
	Map(result interface{}, obj Record) error
}

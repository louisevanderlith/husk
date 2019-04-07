package husk

// Filterer used to filter records while searching.
type Filterer interface {
	// Enables a single central function to cast from husk.Dataer to *<Person>
	Filter(obj Dataer) bool
}

type filter func(obj Dataer) bool

//Filter is the function which casts objects before sending them to the filter.
func (f filter) Filter(obj Dataer) bool {
	return f(obj)
}

// Everything, returns 'true' on all rows
func Everything() filter {
	return func(obj Dataer) bool {
		return true
	}
}

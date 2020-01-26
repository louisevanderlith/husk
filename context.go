package husk

// Ctxer is a special interface which can be used when building extensions
type Ctxer interface {
	//Save calls the save method for every registered table
	Save() error
}
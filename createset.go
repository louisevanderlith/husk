package husk

//CreateSet is a wrapper for a newly created record,
//and any errors that may have caused the creation to fail.
type CreateSet struct {
	Record Recorder
	Error  error
}

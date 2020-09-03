package hsk

//Index is used to manage data and where it's located
type Index interface {
	/// Create new entry in this index that maps value.Key K to value V
	Add(m Meta) (Key, error)
	Set(k Key, v Meta) error
	/// Find an entry by key
	Get(k Key) Meta
	IndexOf(k Key) int
	/// Delete all entries of given key
	Delete(k Key) bool

	GetKeys() []Key
}

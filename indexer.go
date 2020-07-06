package husk

//Indexer is used to manage data and where it's located
type Indexer interface {
	CreateSpace(point *Point) *meta
	/// Create new entry in this index that maps value.Key K to value V
	Insert(v *meta) Key

	/// Find an entry by key
	Get(k Key) *meta
	/// Delete all entries of given key
	Delete(k Key) bool

	Entries() []Key
}

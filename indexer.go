package husk

type Indexer interface {
	CreateSpace(point *Point) *meta
	/// Create new entry in this index that maps value.Key K to value V
	Insert(v *meta)

	/// Find an entry by key
	Get(k *Key) *meta

	/// Delete all entries of given key
	Delete(k *Key) bool

	Items() map[*Key]*meta
}

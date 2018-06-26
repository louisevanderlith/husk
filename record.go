package husk

type Record struct {
	meta *meta
	data Dataer
}

func NewRecord(tableName string, id int64, obj Dataer) Record {
	meta := NewMeta(tableName, id)

	return MakeRecord(meta, obj)
}

func MakeRecord(meta *meta, obj Dataer) Record {
	result := Record{}
	result.meta = meta
	result.data = obj

	return result
}

func (r Record) GetID() int64 {
	return r.meta.ID
}

func (r Record) Meta() *meta {
	return r.meta
}

func (r Record) Data() Dataer {
	return r.data
}

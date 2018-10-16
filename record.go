package husk

type Record struct {
	meta *meta
	data Dataer
}

func MakeRecord(meta *meta, obj Dataer) *Record {
	result := Record{}
	result.meta = meta
	result.data = obj

	return &result
}

func (r Record) GetKey() *Key {
	return r.meta.GetKey()
}

func (r Record) Meta() *meta {
	return r.meta
}

func (r Record) Data() Dataer {
	return r.data
}

func (r *Record) Set(obj Dataer) error {
	valid, err := obj.Valid()

	if err != nil || !valid {
		return err
	}

	r.data = obj

	return nil
}

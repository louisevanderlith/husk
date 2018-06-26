package husk

type Index map[int64]*meta

func LoadIndex(indexName string) *Index {
	result := make(Index)

	err := read(indexName, &result)

	if err != nil {
		panic(err)
	}

	return &result
}

func (m *Index) nextID() int64 {
	var result int64

	for k := range *m {
		if result < k {
			result = k
		}
	}

	return result + 1
}

func (m *Index) getAt(id int64) *meta {
	meta := (*m)[id]

	if meta != nil && meta.Active {
		return meta
	}

	return nil
}

func (m *Index) addMeta(obj *meta) {
	id := obj.ID
	(*m)[id] = obj
}

func (m *Index) dump(tableName string) {
	indexName := getIndexName(tableName)

	err := write(indexName, m)

	if err != nil {
		panic(err)
	}
}

package hsk

type Point interface {
	GetOffset() int64
	GetLength() int64
}

//Point contains the information for finding a value on the Tape
type point struct {
	Offset int64
	Len    int64
}

func NewPoint(offset, length int64) Point {
	return point{Offset: offset, Len: length}
}

func (p point) GetOffset() int64 {
	return p.Offset
}

func (p point) GetLength() int64 {
	return p.Len
}

package husk

type Point struct {
	Offset int64
	Len    int64
}

func NewPoint(offset, length int64) *Point {
	return &Point{Offset: offset, Len: length}
}

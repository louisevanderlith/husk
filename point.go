package husk

//Point contains the information for finding a value on the Tape
type Point struct {
	Offset int64
	Len    int64
}

func newPoint(offset, length int64) *Point {
	return &Point{Offset: offset, Len: length}
}

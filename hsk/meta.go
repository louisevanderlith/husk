package hsk

type Meta interface {
	Disable()
	Update(p *Point)
	GetKey() Key
	IsActive() bool
	Point() *Point
}

type meta struct {
	Key     Key
	Active  bool
	Pointer *Point
}

func NewMeta(key Key, point *Point) Meta {
	return &meta{
		Key:     key,
		Active:  true,
		Pointer: point,
	}
}

func (m *meta) Disable() {
	m.Active = false
}

func (m *meta) Update(p *Point) {
	m.Pointer = p
}

func (m *meta) GetKey() Key {
	return m.Key
}

func (m *meta) IsActive() bool {
	return m.Active
}

func (m *meta) Point() *Point {
	return m.Pointer
}

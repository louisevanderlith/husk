package hsk

type Meta interface {
	IsActive() bool
	Disable()
	Enable()

	Update(p Point)
	Point() Point
}

type meta struct {
	Active  bool
	Pointer Point
}

func NewMeta() Meta {
	return &meta{}
}

func NewMetaWithPoint(p Point) Meta {
	return &meta{
		Active:  true,
		Pointer: p,
	}
}

func (m *meta) Enable() {
	m.Active = true
}
func (m *meta) Disable() {
	m.Active = false
}

func (m *meta) Update(p Point) {
	m.Enable()
	m.Pointer = p
}

func (m *meta) IsActive() bool {
	return m.Active
}

func (m *meta) Point() Point {
	return m.Pointer
}

package husk

import (
	"fmt"
	"time"
)

type Key struct {
	Stamp int64
	ID    int64
}

// CrazyKey is a short-hand for NewKey(-1), returns old date
func CrazyKey() Key {
	old := time.Date(1991, 8, 2, 12, 13, 57, 000, time.UTC)
	return Key{old.UnixNano(), int64(-1)}
}

func NewKey(nextID int64) Key {
	if nextID == -1 {
		panic("rather call CrazyKey")
	}

	timestamp := time.Now().UnixNano()
	return Key{timestamp, nextID}
}

func (k Key) String() string {
	return fmt.Sprintf("%d-%d", k.Stamp, k.ID)
}

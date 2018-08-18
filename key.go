package husk

import (
	"fmt"
	"time"
)

type Key struct {
	Stamp int64
	ID    int64
}

func NewKey(nextID int64) Key {
	timestamp := time.Now().UnixNano()
	return Key{timestamp, nextID}
}

func (k Key) String() string {
	return fmt.Sprintf("%d-%d", k.Stamp, k.ID)
}

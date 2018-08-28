package husk

import (
	"fmt"
	"strconv"
	"strings"
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

func ParseKey(rawKey string) Key {
	parts := strings.Split(rawKey, "-")

	if len(parts) != 2 {
		return CrazyKey()
	}

	stamp, _ := strconv.ParseInt(parts[0], 10, 64)
	id, _ := strconv.ParseInt(parts[1], 10, 64)

	return Key{Stamp: stamp, ID: id}
}

func (k Key) String() string {
	return fmt.Sprintf("%d-%d", k.Stamp, k.ID)
}

func (k Key) GetTimestamp() time.Time {
	return time.Unix(k.Stamp, 0)
}

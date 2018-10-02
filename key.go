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
func CrazyKey() *Key {
	old := time.Date(1991, 8, 2, 12, 13, 57, 000, time.UTC)
	return &Key{old.Unix(), int64(-1)}
}

func NewKey(nextID int64) *Key {
	if nextID == -1 {
		panic("rather call CrazyKey")
	}

	timestamp := time.Now().Unix()
	return &Key{timestamp, nextID}
}

func ParseKey(rawKey string) *Key {
	parts := strings.Split(rawKey, "-")

	if len(parts) != 2 {
		return CrazyKey()
	}

	stamp, _ := strconv.ParseInt(parts[0], 10, 64)
	id, _ := strconv.ParseInt(parts[1], 10, 64)

	return &Key{stamp, id}
}

func (k *Key) ItemNo() int64 {
	return k.ID
}

func (k *Key) Timestamp() int64 {
	return k.Stamp
}

func (k *Key) String() string {
	return fmt.Sprintf("%d-%d", k.Stamp, k.ID)
}

func (k *Key) GetTimestamp() time.Time {
	return time.Unix(k.Stamp, 0)
}

//Compare returns -1 (smaller), 0 (equal), 1 (larger)
func (k *Key) Compare(k2 *Key) int8 {
	//Stamps are checked before ID
	if k.Stamp < k2.Timestamp() {
		return -1
	}

	if k.ID > k2.ItemNo() {
		return 1
	}

	//Stamps are Equal
	if k.ID < k2.ItemNo() {
		return -1
	}

	if k.ID > k2.ItemNo() {
		return 1
	}

	return 0
}

package husk

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

//Key is the Primary Key for Husk Indexes.
type Key struct {
	//Stamp is the EPOCH Creation Time
	Stamp int64
	//ID increments with dulpicate Stamps
	ID int64
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

// ParseKey tries to parse EPOCH-00 Keys.
func ParseKey(rawKey string) (*Key, error) {
	parts := strings.Split(rawKey, "-")

	if len(parts) != 2 {
		return nil, errors.New("key not valid format")
	}

	stamp, err := strconv.ParseInt(parts[0], 10, 64)

	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(parts[1], 10, 64)

	if err != nil {
		return nil, err
	}

	return &Key{stamp, id}, nil
}

//String returns the string representation for a Key, also makes is easier to parse.
func (k *Key) String() string {
	return fmt.Sprintf("%d-%d", k.Stamp, k.ID)
}

//Timestamp gets the Stamp value of the Key
func (k *Key) GetTimestamp() time.Time {
	return time.Unix(k.Stamp, 0)
}

//Compare returns -1 (smaller), 0 (equal), 1 (larger)
func (k *Key) Compare(k2 *Key) int8 {
	//Stamps are checked before ID
	if k.Stamp < k2.Stamp {
		return -1
	}

	if k.ID > k2.ID {
		return 1
	}

	//Stamps are Equal
	if k.ID < k2.ID {
		return -1
	}

	if k.ID > k2.ID {
		return 1
	}

	return 0
}

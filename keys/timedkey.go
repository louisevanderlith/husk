package keys

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/louisevanderlith/husk/hsk"
	"strconv"
	"strings"
	"time"
)

//TimeKey is the Primary Key for Husk Indexes. It orders records by Timestamp
type TimeKey struct {
	//Stamp is the EPOCH Creation Time
	Stamp int64
	//ID increments with duplicate Stamps
	ID int64
}

// CrazyKey is a short-hand for NewKey(-1), returns old date
func CrazyKey() hsk.Key {
	old := time.Date(1991, 8, 2, 12, 13, 57, 000, time.UTC)
	return NewKeyWithTime(old.Unix(), int64(-1))
}

//NewKey creates a new key with the current time as the Stamp
func NewKey(nextID int64) hsk.Key {
	if nextID == -1 {
		panic("rather call CrazyKey")
	}

	timestamp := time.Now().Unix()
	return NewKeyWithTime(timestamp, nextID)
}

func NewKeyWithTime(timestamp, nextID int64) hsk.Key {
	return &TimeKey{timestamp, nextID}
}

// ParseKey tries to parse EPOCH`00 Keys.
func ParseKey(rawKey string) (hsk.Key, error) {
	if len(rawKey) < 3 {
		return nil, errors.New("rawKey too short")
	}

	if !strings.Contains(rawKey, "`") {
		return nil, errors.New("key not valid format")
	}

	dotIndx := strings.Index(rawKey, "`")
	stamp, err := strconv.ParseInt(rawKey[:dotIndx], 10, 64)

	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseInt(rawKey[dotIndx+1:], 10, 64)

	if err != nil {
		return nil, err
	}

	return &TimeKey{stamp, id}, nil
}

//String returns the string representation for a Key, also makes is easier to parse.
func (k *TimeKey) String() string {
	return fmt.Sprintf("%d`%d", k.Stamp, k.ID)
}

//GetTimestamp returns the creation time of the record
func (k *TimeKey) GetTimestamp() time.Time {
	return time.Unix(k.Stamp, 0)
}

//Compare returns -1 (smaller), 0 (equal), 1 (larger)
func (k *TimeKey) Compare(obj interface{}) int8 {
	k2 := obj.(*TimeKey)

	//Stamps are checked before ID
	if k.Stamp < k2.Stamp {
		return -1
	}

	if k.Stamp > k2.Stamp {
		return 1
	}

	if k.ID < k2.ID {
		return -1
	}

	if k.ID > k2.ID {
		return 1
	}

	return 0
}

//MarshalJSON will return a Key as {stamp}`{key}
func (k *TimeKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

//UnmarshalJSON will return {stamp}`{key} as a Key
func (k *TimeKey) UnmarshalJSON(b []byte) error {
	stripEsc := strings.Replace(string(b), "\"", "", -1)
	tmpK, err := ParseKey(stripEsc)

	if err != nil {
		return err
	}

	tmKey, ok := tmpK.(*TimeKey)

	if !ok {
		return errors.New("invalid key type")
	}

	k.Stamp = tmKey.Stamp
	k.ID = tmKey.ID

	return nil
}

package hsk

import (
	"encoding/json"
	"fmt"
	"github.com/louisevanderlith/husk/collections"
	"time"
)

type Key interface {
	collections.Comparable
	fmt.Stringer
	json.Marshaler
	json.Unmarshaler
	GetTimestamp() time.Time
}

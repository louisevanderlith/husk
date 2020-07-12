package collections

import (
	"github.com/louisevanderlith/husk/hsk"
)

//Comparable has the Compare Function
type Comparable interface {
	//Compare returns -1 (smaller), 0 (equal), 1 (larger)
	Compare(k hsk.Key) int8
}

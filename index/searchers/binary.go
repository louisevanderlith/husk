package searchers

import (
	"github.com/louisevanderlith/husk/hsk"
	"sort"
)

//Return the index of a key using binary search
func IndexOf(keys []hsk.Key, key hsk.Key) int {
	idx := sort.Search(len(keys), func(n int) bool {
		curr := keys[n]
		//Smaller or Equals, since husk is ordered by Created Date desc
		comp := curr.Compare(key)
		return comp <= 0
	})

	if idx < len(keys) && keys[idx].Compare(key) == 0 {
		return idx
	}

	return -1
}

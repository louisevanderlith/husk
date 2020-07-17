package hybrid

import (
	"github.com/louisevanderlith/husk/hsk"
	"github.com/louisevanderlith/husk/storers"
)

func NewTable(obj hsk.Dataer) storers.Table {
	return storers.NewTable(obj, newStore())
}

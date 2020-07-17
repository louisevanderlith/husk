package tape

import (
	"github.com/louisevanderlith/husk/db"
	"testing"
)

func TestNewTable(t *testing.T) {
	tbl := NewTable(db.Event{})

	if tbl.Name() != "Event" {
		t.Error("Invalid Name; expected Event got", tbl.Name())
	}
}

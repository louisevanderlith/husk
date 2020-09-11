package records

import (
	"encoding/json"
	"github.com/louisevanderlith/husk/keys"
	"log"
	"testing"
)

func TestRecord_MarshalJSON(t *testing.T) {
	in := MakeRecord(keys.NewKeyWithTime(1599470402, 5), &alpha{"XF"})
	data, err := json.Marshal(in)

	if err != nil {
		t.Error(err)
		return
	}

	exp := "{\"Key\":\"1599470402`5\",\"Value\":{\"Char\":\"XF\"}}"

	if string(data) != exp {
		t.Error("Expected", exp, "Got", string(data))
	}
}

func TestRecord_UnmarshalJSON(t *testing.T) {
	in := "{\"Key\":\"1599470402`5\",\"Value\":{\"Char\":\"XF\"}}"
	act := NewRecord(&alpha{})
	err := json.Unmarshal([]byte(in), &act)

	if err != nil {
		t.Error(err)
		return
	}

	exp := MakeRecord(keys.NewKeyWithTime(1599470402, 5), &alpha{"XF"})

	if act.GetKey().Compare(exp.GetKey()) != 0 {
		t.Error("Expected Key", exp, "Got", act)
	}

	log.Println("Value is", act.GetValue().(*alpha))
	alp := act.GetValue().(*alpha)

	if alp.Char != "XF" {
		t.Error("Unexpected Value", alp.Char)
	}
}

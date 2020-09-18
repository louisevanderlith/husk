package records

import (
	"encoding/json"
	"github.com/louisevanderlith/husk/keys"
	"testing"
)

func TestPage_Marshal_UnMarshalJSON(t *testing.T) {
	in := NewRecordPage(2, 5)
	in.Add(MakeRecord(keys.NewKey(6), &alpha{"F"}))
	in.Add(MakeRecord(keys.NewKey(7), &alpha{"G"}))
	in.Add(MakeRecord(keys.NewKey(8), &alpha{"H"}))
	in.Add(MakeRecord(keys.NewKey(9), &alpha{"I"}))
	in.Add(MakeRecord(keys.NewKey(10), &alpha{"J"}))

	act, err := json.Marshal(in)

	if err != nil {
		t.Error(err)
		return
	}

	out := NewResultPage(alpha{})
	err = json.Unmarshal(act, out)

	if err != nil {
		t.Error(err)
		return
	}

	if out.Count() != in.Count() {
		t.Error("Expected", out.Count(), "Got", in.Count())
	}
}

func TestPage_MarshalJSON(t *testing.T) {
	in := NewRecordPage(2, 5)

	act, err := json.Marshal(in)

	if err != nil {
		t.Error(err)
		return
	}

	exp := `{"Records":null,"Number":2,"Size":5,"HasMore":false}`

	if string(act) != exp {
		t.Error("Expected", exp, "Got", string(act))
	}
}

func TestPage_UnMarshalJSON(t *testing.T) {
	act := NewResultPage(alpha{})
	in := "{\"Records\":[{\"Key\":\"1599470402`4\",\"Value\":{\"Char\":\"C\"}},{\"Key\":\"1599470402`5\",\"Value\":{\"Char\":\"D\"}}],\"Number\":2,\"Size\":5,\"Limit\":10}"
	err := json.Unmarshal([]byte(in), act)

	if err != nil {
		t.Error(err)
		return
	}

	exp := NewRecordPage(2, 2)
	exp.Add(MakeRecord(keys.NewKeyWithTime(1599470402, 4), &alpha{"C"}))
	exp.Add(MakeRecord(keys.NewKeyWithTime(1599470402, 5), &alpha{"D"}))

	if act.Count() != exp.Count() {
		t.Error("Count Expected", exp.Count(), "Got", act.Count())
	}
}

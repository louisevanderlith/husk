package keys

import (
	"encoding/json"
	"testing"
	"time"
)

func TestKey_CanParse(t *testing.T) {
	k := CrazyKey().String()
	prs, err := ParseKey(k)

	if err != nil {
		t.Error(err)
	}

	if prs.String() != k {
		t.Errorf("Expected %s, got %+v.", k, prs)
	}
}

func TestKey_CanParseEmpty(t *testing.T) {
	k := "0`0"
	prs, err := ParseKey(k)

	if err != nil {
		t.Error(err)
		return
	}

	if prs.String() == CrazyKey().String() {
		t.Errorf("Expected %v, got %v", prs, CrazyKey())
	}
}

func TestKey_ToJSON(t *testing.T) {
	k := CrazyKey()

	expected, _ := k.MarshalJSON()
	actual, err := json.Marshal(k)

	if err != nil {
		t.Error(err)
	}

	if string(actual) != string(expected) {
		t.Errorf("expected %s, got %s", string(expected), string(actual))
	}
}

func TestKey_FromJSON(t *testing.T) {
	expected := NewKey(2)
	input, _ := expected.MarshalJSON()

	actual := CrazyKey()
	err := json.Unmarshal(input, &actual)

	if err != nil {
		t.Error(err)
	}

	if actual.String() != expected.String() {
		t.Errorf("expected %v, got %v", expected, actual)
	}
}

var (
	y67 = time.Date(1967, 06, 06, 13, 15, 16, 000, time.UTC).UnixNano()
	y99 = time.Date(1999, 06, 06, 13, 15, 16, 000, time.UTC).UnixNano()
)

func TestLess_SameDate_SameID_AreEqual(t *testing.T) {
	iK := NewKeyWithTime(y67, 22)
	jK := NewKeyWithTime(y67, 22)

	if iK.Compare(jK) != 0 {
		t.Error("Expected Equal", iK.Compare(jK))
	}
}

func TestLess_SameDate_SmallerID_iK_Smaller(t *testing.T) {
	iK := NewKeyWithTime(y67, 8)
	jK := NewKeyWithTime(y67, 22)

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_SameDate_LargerID_iK_Larger(t *testing.T) {
	iK := NewKeyWithTime(y67, 22)
	jK := NewKeyWithTime(y67, 8)

	if iK.Compare(jK) != 1 {
		t.Error("Expected iK Larger", iK.Compare(jK))
	}
}

func TestLess_SmallerDate_SameID_iK_Smaller(t *testing.T) {
	iK := NewKeyWithTime(y67, 22)
	jK := NewKeyWithTime(y99, 22)

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_SmallerDate_SmallerID_iK_Smaller(t *testing.T) {
	iK := NewKeyWithTime(y67, 8)
	jK := NewKeyWithTime(y99, 22)

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_SmallerDate_LargerID_iK_Smaller(t *testing.T) {
	iK := NewKeyWithTime(y67, 22)
	jK := NewKeyWithTime(y99, 8)

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_LargerDate_SameID_iK_Larger(t *testing.T) {
	iK := NewKeyWithTime(y99, 22)
	jK := NewKeyWithTime(y67, 22)

	if iK.Compare(jK) != 1 {
		t.Error("Expected iK Larger", iK.Compare(jK))
	}
}

func TestLess_LargerDate_SmallerID_iK_Larger(t *testing.T) {
	iK := NewKeyWithTime(y99, 8)
	jK := NewKeyWithTime(y67, 22)

	if iK.Compare(jK) != 1 {
		t.Error("Expected iK Larger", iK.Compare(jK))
	}
}

func TestLess_LargerDate_LargerID_iK_Smaller(t *testing.T) {
	iK := NewKeyWithTime(y67, 22)
	jK := NewKeyWithTime(y99, 89)

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller")
	}
}

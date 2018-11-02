package husk

import (
	"testing"
	"time"
)

var (
	y67 = time.Date(1967, 06, 06, 13, 15, 16, 000, time.UTC).UnixNano()
	y99 = time.Date(1999, 06, 06, 13, 15, 16, 000, time.UTC).UnixNano()
)

func TestLess_SameDate_SameID_AreEqual(t *testing.T) {
	iK := &Key{y67, 22}
	jK := &Key{y67, 22}

	if iK.Compare(jK) != 0 {
		t.Error("Expected Equal", iK.Compare(jK))
	}
}

func TestLess_SameDate_SmallerID_iK_Smaller(t *testing.T) {
	iK := &Key{y67, 8}
	jK := &Key{y67, 22}

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_SameDate_LargerID_iK_Larger(t *testing.T) {
	iK := &Key{y67, 22}
	jK := &Key{y67, 8}

	if iK.Compare(jK) != 1 {
		t.Error("Expected iK Larger", iK.Compare(jK))
	}
}

func TestLess_SmallerDate_SameID_iK_Smaller(t *testing.T) {
	iK := &Key{y67, 22}
	jK := &Key{y99, 22}

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_SmallerDate_SmallerID_iK_Smaller(t *testing.T) {
	iK := &Key{y67, 8}
	jK := &Key{y99, 22}

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_SmallerDate_LargerID_iK_Smaller(t *testing.T) {
	iK := &Key{y67, 22}
	jK := &Key{y99, 8}

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller", iK.Compare(jK))
	}
}

func TestLess_LargerDate_SameID_iK_Larger(t *testing.T) {
	iK := &Key{y99, 22}
	jK := &Key{y67, 22}

	if iK.Compare(jK) != 1 {
		t.Error("Expected iK Larger", iK.Compare(jK))
	}
}

func TestLess_LargerDate_SmallerID_iK_Larger(t *testing.T) {
	iK := &Key{y99, 8}
	jK := &Key{y67, 22}

	if iK.Compare(jK) != 1 {
		t.Error("Expected iK Larger", iK.Compare(jK))
	}
}

func TestLess_LargerDate_LargerID_iK_Smaller(t *testing.T) {
	iK := &Key{y67, 22}
	jK := &Key{y99, 89}

	if iK.Compare(jK) != -1 {
		t.Error("Expected iK Smaller")
	}
}

func TestLoadIndex_AllDataPresent(t *testing.T) {
	indxName := getIndexName("Person")
	indx := loadIndex(indxName)

	if len(indx.Items()) == 0 {
		t.Error("No data")
	}

	t.Fail()
}

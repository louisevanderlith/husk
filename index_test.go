package husk

import (
	"testing"
	"time"
)

var (
	y67 = time.Date(1967, 06, 06, 13, 15, 16, 000, time.UTC).UnixNano()
	y99 = time.Date(1999, 06, 06, 13, 15, 16, 000, time.UTC).UnixNano()
)

func TestZero(t *testing.T) {

}

func TestLess_SameDate_SameID_MustFalse(t *testing.T) {
	iK := Key{y67, 22}
	jK := Key{y67, 22}

	if isLess(iK, jK) {
		t.Error("Expected False")
	}
}

func TestLess_SameDate_SmallerID_MustTrue(t *testing.T) {
	iK := Key{y67, 8}
	jK := Key{y67, 22}

	if !isLess(iK, jK) {
		t.Error("Expected True")
	}
}

func TestLess_SameDate_LargerID_MustFalse(t *testing.T) {
	iK := Key{y67, 22}
	jK := Key{y67, 8}

	if isLess(iK, jK) {
		t.Error("Expected False")
	}
}

func TestLess_SmallerDate_SameID_MustTrue(t *testing.T) {
	iK := Key{y67, 22}
	jK := Key{y99, 22}

	if !isLess(iK, jK) {
		t.Error("Expected True")
	}
}

func TestLess_SmallerDate_SmallerID_MustTrue(t *testing.T) {
	iK := Key{y67, 8}
	jK := Key{y99, 22}

	if !isLess(iK, jK) {
		t.Error("Expected True")
	}
}

func TestLess_SmallerDate_LargerID_MustTrue(t *testing.T) {
	iK := Key{y67, 22}
	jK := Key{y99, 8}

	if !isLess(iK, jK) {
		t.Error("Expected True")
	}
}

func TestLess_LargerDate_SameID_MustFalse(t *testing.T) {
	iK := Key{y99, 22}
	jK := Key{y67, 22}

	if isLess(iK, jK) {
		t.Error("Expected False")
	}
}

func TestLess_LargerDate_SmallerID_MustFalse(t *testing.T) {
	iK := Key{y99, 8}
	jK := Key{y67, 22}

	if isLess(iK, jK) {
		t.Error("Expected False")
	}
}

func TestLess_LargerDate_LargerID_MustFalse(t *testing.T) {
	iK := Key{y67, 22}
	jK := Key{y67, 22}

	if isLess(iK, jK) {
		t.Error("Expected False")
	}
}

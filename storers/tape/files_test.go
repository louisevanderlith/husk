package tape

import (
	"testing"
)

func TestGetRecordName_Correct(t *testing.T) {
	actual := getDataPath("Person")
	expected := "db/Person.Data.husk"

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

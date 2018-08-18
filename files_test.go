package husk

import (
	"fmt"
	"testing"
)

func TestGetRecordName_Correct(t *testing.T) {
	k := NewKey(9)
	actual := getRecordName("Person", k)
	expected := fmt.Sprintf("db/Person.%s.husk", k)

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestGetIndexName_Correct(t *testing.T) {
	expected := "db/Person.index.husk"
	actual := getIndexName("Person")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestCleanTableName_Correct(t *testing.T) {
	expected := "Person"
	actual := cleanTableName("Person.index.husk")

	if actual != expected {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

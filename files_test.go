package husk

import "testing"

func TestGetRecordName_Correct(t *testing.T) {
	expected := "db/Person.9.husk"
	actual := getRecordName("Person", 9)

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

package ozone

import (
	"os"
	"testing"
)

func TestSaveAndLoadXMLFile(t *testing.T) {
	// Creating and initializing the database
	db := New(2)
	defer db.Close()

	// Adding items to the database
	err := db.Set("key1", "value1")
	if err != nil {
		t.Fatalf("Error setting value: %v", err)
	}

	err = db.Set("key2", "value2")
	if err != nil {
		t.Fatalf("Error setting value: %v", err)
	}

	// Saving database to a compressed XML file
	err = db.Save("test_db.oz")
	if err != nil {
		t.Fatalf("Error saving to XML file: %v", err)
	}

	// Creating and initializing a new database for loading
	newDB := New(2)
	defer newDB.Close()

	// Loading database from a compressed XML file
	err = newDB.Load("test_db.oz")
	if err != nil {
		t.Fatalf("Error loading from XML file: %v", err)
	}

	// Cleaning up the saved file
	err = os.Remove("test_db.oz")
	if err != nil {
		t.Fatalf("Error cleaning up test file: %v", err)
	}

	// Verifying the loaded data
	dataItem, err := newDB.Get("key1")
	if err != nil {
		t.Fatalf("Error getting value: %v", err)
	}
	if dataItem.Value != "value1" {
		t.Fatalf("Expected value1, but got %v", dataItem.Value)
	}

	dataItem, err = newDB.Get("key2")
	if err != nil {
		t.Fatalf("Error getting value: %v", err)
	}
	if dataItem.Value != "value2" {
		t.Fatalf("Expected value2, but got %v", dataItem.Value)
	}
}

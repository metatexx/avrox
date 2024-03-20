package rawdate_test

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/metatexx/avrox/rawdate"
)

// TestRawDateInSQLite3DATE tests the RawDate custom type against an SQLite in-memory database
// This is using the DATE type in the DB
func TestRawDateInSQLite3DATE(t *testing.T) {
	// Connect to the SQLite in-memory database.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create a table with a Date column.
	_, err = db.Exec("CREATE TABLE test_table (Date DATE)")
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// Prepare a date for testing.
	testDate := rawdate.Today()

	// Insert the test date into the table.
	_, err = db.Exec("INSERT INTO test_table (Date) VALUES (?)", testDate)
	if err != nil {
		t.Fatalf("Error inserting data: %v", err)
	}

	// Retrieve the date from the table.
	var retrievedDate rawdate.RawDate
	err = db.QueryRow("SELECT Date FROM test_table LIMIT 1").Scan(&retrievedDate)
	if err != nil {
		t.Fatalf("Error retrieving data: %v", err)
	}

	// Compare the inserted and retrieved dates.
	if retrievedDate != testDate {
		t.Errorf("Retrieved date %v does not match inserted date %v", retrievedDate, testDate)
	}
}

// TestRawDateInSQLite3TEXT tests the RawDate custom type against an SQLite in-memory database
// This is using the TEXT type in the DB (to proof a point actually)
func TestRawDateInSQLite3TEXT(t *testing.T) {
	// Connect to the SQLite in-memory database.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create a table with a Date column.
	_, err = db.Exec("CREATE TABLE test_table (Date TEXT)")
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// Prepare a date for testing.
	testDate := rawdate.Today()

	// Insert the test date into the table.
	_, err = db.Exec("INSERT INTO test_table (Date) VALUES (?)", testDate)
	if err != nil {
		t.Fatalf("Error inserting data: %v", err)
	}

	// Retrieve the date from the table. This will parse the TEXT back to the RawDate
	var retrievedDate rawdate.RawDate
	err = db.QueryRow("SELECT Date FROM test_table LIMIT 1").Scan(&retrievedDate)
	if err != nil {
		t.Fatalf("Error retrieving data: %v", err)
	}

	// Compare the inserted and retrieved dates.
	if retrievedDate != testDate {
		t.Errorf("Retrieved date %v does not match inserted date %v", retrievedDate, testDate)
	}
}

// TestRawDateInSQLite3TEXTTime tests the RawDate custom type against an SQLite in-memory database
// This is using the TEXT type in the DB and parse creates a time.Time from the RawDate
func TestRawDateInSQLite3TEXTTime(t *testing.T) {
	// Connect to the SQLite in-memory database.
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	// Create a table with a Date column but make it the TEXT type
	_, err = db.Exec("CREATE TABLE test_table (Date TEXT)")
	if err != nil {
		t.Fatalf("Error creating table: %v", err)
	}

	// Prepare a date for testing.
	testDate := rawdate.Today()

	rawdate.SQLValueTime = true
	// Insert the test date into the table.
	_, err = db.Exec("INSERT INTO test_table (Date) VALUES (?)", testDate)
	rawdate.SQLValueTime = false
	if err != nil {
		t.Fatalf("Error inserting data: %v", err)
	}

	// Retrieve the date from the table as a string.
	var retrievedDate string
	err = db.QueryRow("SELECT Date FROM test_table LIMIT 1").Scan(&retrievedDate)
	if err != nil {
		t.Fatalf("Error retrieving data: %v", err)
	}
	// Compare the inserted and retrieved dates.
	if retrievedDate != testDate.Format("2006-01-02 00:00:00+00:00") {
		t.Errorf("Retrieved date %v does not match inserted date %v", retrievedDate, testDate)
	}
}

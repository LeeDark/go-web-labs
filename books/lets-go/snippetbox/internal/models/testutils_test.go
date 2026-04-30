package models

import (
	"database/sql"
	"os"
	"testing"
)

func newTestDB(t *testing.T) *sql.DB {
	// Establish a sql.DB connection pool for our test database
	db, err := sql.Open("mysql", "test_web:Pass_Test1@/test_snippetbox?parseTime=true&multiStatements=true")
	if err != nil {
		t.Fatal(err)
	}

	// Read the setup SQL script
	script, err := os.ReadFile("./testdata/setup.sql")
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	// Execute statements from SQL script
	_, err = db.Exec(string(script))
	if err != nil {
		db.Close()
		t.Fatal(err)
	}

	t.Cleanup(func() {
		defer db.Close()

		// Read the teardown SQL script
		script, err := os.ReadFile("./testdata/teardown.sql")
		if err != nil {
			t.Fatal(err)
		}

		// Execute statements from SQL script
		_, err = db.Exec(string(script))
		if err != nil {
			t.Fatal(err)
		}
	})

	return db
}

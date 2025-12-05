package store

import (
	"database/sql"
	"flag"
	"fmt"
	"os"
	"testing"

	"h3jfc/shed/lib/sqlite3"
)

var database *sql.DB

const dbname = "file::memory:?cache=shared&_journal_mode=WAL&_busy_timeout=10000"

func TestMain(m *testing.M) {
	var err error

	flag.Parse()

	if testing.Short() {
		os.Exit(m.Run())
	}

	database, err = sql.Open("sqlite3", dbname)
	if err != nil {
		panic(fmt.Sprintf("failed to open database: %v", err))
	}

	database.SetMaxOpenConns(1)

	if err := sqlite3.MigrateDB(database); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	code := m.Run()

	if err := database.Close(); err != nil {
		panic(fmt.Sprintf("failed to close database: %v", err))
	}

	os.Exit(code)
}

func prepTx(t *testing.T) *sql.Tx {
	t.Helper()

	tx, err := database.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	t.Cleanup(func() {
		if err := tx.Rollback(); err != nil {
			t.Fatalf("failed to rollback transaction: %v", err)
		}
	})

	return tx
}

func prepNewStore(t *testing.T) *Store {
	t.Helper()

	tx := prepTx(t)

	return NewStore(tx)
}

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"h3jfc/shed/cmd"
	sheddb "h3jfc/shed/db"
	"h3jfc/shed/lib/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

var (
	targetVersion     = 1
	ErrDirtyMigration = errors.New("migration is dirty, intervention is needed")
)

func main() {
	err := migrateSqlite3()
	checkErr(err)
	cmd.Execute()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func migrateSqlite3() error {
	// Ensure .shed directory exists
	shedDir, err := config.GetShedDir()
	if err != nil {
		return fmt.Errorf("failed to get shed directory: %w", err)
	}

	if err := os.MkdirAll(shedDir, 0755); err != nil {
		return fmt.Errorf("failed to create .shed directory: %w", err)
	}

	key := "my_secret_key" // In a real application, use a secure method to manage encryption keys.
	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return fmt.Errorf("failed to get database path: %w", err)
	}

	dbname := fmt.Sprintf("file:%s?_key=%s&_cipher_page_size=%d", dbPath, key, 4096)

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}

	defer func() {
		if err := db.Close(); err != nil {
			return
		}
	}()

	sourceDriver, err := iofs.New(sheddb.Migrations, "migrations")
	if err != nil {
		return fmt.Errorf("could not access pgx migrations FS: %w", err)
	}

	dbDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		return fmt.Errorf("could not create sqlite3 driver: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",       // source name
		sourceDriver, // source driver
		"sqlite3",    // database name
		dbDriver,     // database driver
	)
	if err != nil {
		return err
	}

	currentVersion, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get current migration version: %w", err)
	}

	if dirty {
		return ErrDirtyMigration
	}

	if targetVersion > int(currentVersion) {
		err := m.Migrate(uint(targetVersion))
		if err != nil {
			return fmt.Errorf("could not migrate to version %v: %w", targetVersion, err)
		}
	}

	return nil
}

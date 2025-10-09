package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"h3jfc/shed/cmd"
	sheddb "h3jfc/shed/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3"
)

var (
	targetVersion = 1
	// ErrDirtyMigration error when a migration has errored and needs to be fixed.
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
	dir, err := os.UserHomeDir()
	checkErr(err)

	db, err := sql.Open("sqlite3", filepath.Join(dir, "sqlite3.db"))
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

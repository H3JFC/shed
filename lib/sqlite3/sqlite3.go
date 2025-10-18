package sqlite3

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	msqlite3 "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Register file source driver
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/mattn/go-sqlite3" // Register SQLite3 driver with SQLCipher support

	sheddb "h3jfc/shed/db"
	"h3jfc/shed/internal/config"
)

const (
	defaultTargetVersion  = 1
	defaultDirPerms       = 0o755
	defaultCipherPageSize = 4096
)

var ErrDirtyMigration = errors.New("migration is dirty, intervention is needed")

func Migrate(encryptionKey string) error {
	db, err := SetupDatabase(encryptionKey)
	if err != nil {
		return err
	}
	defer closeDatabase(db)

	m, err := createMigrator(db)
	if err != nil {
		return err
	}

	return runMigrations(m)
}

func SetupDatabase(encryptionKey string) (*sql.DB, error) {
	if err := ensureShedDirectoryExists(); err != nil {
		return nil, err
	}

	dbPath, err := config.GetDatabasePath()
	if err != nil {
		return nil, fmt.Errorf("failed to get database path: %w", err)
	}

	dbname := fmt.Sprintf("file:%s?_key=%s&_cipher_page_size=%d", dbPath, encryptionKey, defaultCipherPageSize)

	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func ensureShedDirectoryExists() error {
	shedDir, err := config.GetShedDir()
	if err != nil {
		return fmt.Errorf("failed to get shed directory: %w", err)
	}

	if err := os.MkdirAll(shedDir, defaultDirPerms); err != nil {
		return fmt.Errorf("failed to create .shed directory: %w", err)
	}

	return nil
}

func closeDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		// Log error but don't return it in defer
		return
	}
}

func createMigrator(db *sql.DB) (*migrate.Migrate, error) {
	sourceDriver, err := iofs.New(sheddb.Migrations, "migrations")
	if err != nil {
		return nil, fmt.Errorf("could not access migrations FS: %w", err)
	}

	dbDriver, err := msqlite3.WithInstance(db, &msqlite3.Config{})
	if err != nil {
		return nil, fmt.Errorf("could not create sqlite3 driver: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",       // source name
		sourceDriver, // source driver
		"sqlite3",    // database name
		dbDriver,     // database driver
	)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func runMigrations(m *migrate.Migrate) error {
	currentVersion, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get current migration version: %w", err)
	}

	if dirty {
		return ErrDirtyMigration
	}

	// Safe conversion: defaultTargetVersion is a const, currentVersion comes from migrate
	if uint(defaultTargetVersion) > currentVersion {
		if err := m.Migrate(uint(defaultTargetVersion)); err != nil {
			return fmt.Errorf("could not migrate to version %v: %w", defaultTargetVersion, err)
		}
	}

	return nil
}

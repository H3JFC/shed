package config

import (
	"bytes"
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"

	"h3jfc/shed/internal/logger"
)

var buf bytes.Buffer

func init() {
	logger.SetWriter(&buf)
	logger.New(logger.ModeVerbose)
}

func TestFindDir_WithEnvironmentVariable(t *testing.T) {
	expectedPath := createValidPathCWD(t)
	t.Setenv("SHED_DIR", expectedPath)

	// Execute
	result, err := FindDir()
	if err != nil {
		t.Fatalf("Expected FindDir() to succeed, but got error: %v", err)
	}

	// Assert
	if result != expectedPath {
		t.Errorf("Expected FindDir() to return %q, but got %q", expectedPath, result)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Found SHED_DIR environment variable")) {
		t.Errorf("Expected debug log about SHED_DIR environment variable, but it was not found")
	}
}

func TestFindDir_NoEnv_DefaultPathFound(t *testing.T) {
	t.Cleanup(clearBuffer)
	t.Setenv("SHED_DIR", "")

	for _, exp := range DefaultConfigPaths { // nolint:paralleltest
		t.Run("test_dir_"+exp, func(t *testing.T) {
			createValidPath(t, exp)

			// Execute
			result, err := FindDir()
			if err != nil {
				t.Fatalf("Expected FindDir() to succeed, but got error: %v", err)
			}

			// Assert
			if result != exp {
				t.Errorf("Expected FindDir() to return %s, but got %q", exp, result)
			}

			if !bytes.Contains(buf.Bytes(), []byte("SHED_DIR environment variable not set")) {
				t.Errorf("Expected debug log about SHED_DIR environment variable not set, but it was not found")
			}
		})
	}
}

func TestFindDir_NothingFound(t *testing.T) {
	t.Cleanup(clearBuffer)
	t.Setenv("SHED_DIR", "")

	// Execute
	result, err := FindDir()
	if err == nil {
		t.Fatalf("Expected FindDir() to return error, but got none")
	}

	if !errors.Is(err, ErrNoPathFound) {
		t.Fatalf("Expected FindDir() to return ErrNoPathFound, but got: %v", err)
	}

	// Assert
	if result != "" {
		t.Errorf("Expected FindDir() to return empty string, but got %q", result)
	}

	if !bytes.Contains(buf.Bytes(), []byte("SHED_DIR environment variable not set")) {
		t.Errorf("Expected debug log about SHED_DIR environment variable not set, but it was not found")
	}

	if !bytes.Contains(buf.Bytes(), []byte("No existing shed path found.")) {
		t.Errorf("Expected debug log about no existing shed path found, but it was not found")
	}
}

func TestShedPath_WithEnvironmentVariable(t *testing.T) {
	expectedPath := createValidPathCWD(t)
	t.Setenv("SHED_DIR", expectedPath)

	// Execute
	result := shedDir()

	// Assert
	if result != expectedPath {
		t.Errorf("Expected shedPath() to return %q, but got %q", expectedPath, result)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Found SHED_DIR environment variable")) {
		t.Errorf("Expected debug log about SHED_DIR environment variable, but it was not found")
	}
}

func TestShedPath_WithoutEnvironmentVariable(t *testing.T) {
	t.Cleanup(clearBuffer)
	t.Setenv("SHED_DIR", "")

	// Execute
	result := shedDir()

	// Assert
	if result != "" {
		t.Errorf("Expected shedPath() to return empty string, but got %q", result)
	}

	if !bytes.Contains(buf.Bytes(), []byte("SHED_DIR environment variable not set")) {
		t.Errorf("Expected debug log about SHED_DIR environment variable not set, but it was not found")
	}
}

func TestValidatePath_ValidPath(t *testing.T) {
	t.Parallel()

	tmpDir := t.TempDir()

	setupValidPath(t, tmpDir)

	if !validatePath(tmpDir) {
		t.Error("validatePath returned false for valid path")
	}
}

func TestValidatePath_PathDoesNotExist(t *testing.T) {
	t.Parallel()

	nonExistentPath := "/path/that/does/not/exist/xyz123"

	if validatePath(nonExistentPath) {
		t.Error("validatePath returned true for non-existent path")
	}
}

func TestValidatePath_PathIsFile(t *testing.T) {
	t.Parallel()
	// Create a temporary file (not a directory)
	tmpFile, err := os.CreateTemp("", "test-file-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	t.Cleanup(func() {
		os.Remove(tmpFile.Name())
	})

	if validatePath(tmpFile.Name()) {
		t.Error("validatePath returned true for file instead of directory")
	}
}

func TestValidatePath_MissingDatabaseFile(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "shed-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mkdir(t, tmpDir)

	configContent := `[shed-db]
password = "test_password"
`
	writeConfigFile(t, tmpDir, configContent)

	if validatePath(tmpDir) {
		t.Error("validatePath returned true when database file is missing")
	}
}

func TestValidatePath_MissingConfigFile(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "shed-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mkdir(t, tmpDir)
	writeDatabaseFile(t, tmpDir)

	if validatePath(tmpDir) {
		t.Error("validatePath returned true when config file is missing")
	}
}

func TestValidatePath_InvalidConfigTOML(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "shed-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mkdir(t, tmpDir)
	writeDatabaseFile(t, tmpDir)

	// Create invalid TOML syntax
	configPath := filepath.Join(tmpDir, "config.toml")

	invalidContent := `[shed-db
this is not valid TOML syntax
`
	if err := os.WriteFile(configPath, []byte(invalidContent), 0o644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	if validatePath(tmpDir) {
		t.Error("validatePath returned true for invalid TOML syntax")
	}
}

func TestValidatePath_MissingShedDbSection(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "shed-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	mkdir(t, tmpDir)
	writeDatabaseFile(t, tmpDir)

	// Create config without [shed-db] section
	configPath := filepath.Join(tmpDir, "config.toml")

	configContent := `[other_section]
key = "value"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	if validatePath(tmpDir) {
		t.Error("validatePath returned true when [shed-db] section is missing")
	}
}

func TestValidatePath_MissingPasswordField(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "shed-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	configContent := `[shed-db]
other_field = "value"
`

	mkdir(t, tmpDir)
	writeConfigFile(t, tmpDir, configContent)
	writeDatabaseFile(t, tmpDir)

	if validatePath(tmpDir) {
		t.Error("validatePath returned true when password field is missing")
	}
}

func TestValidatePath_EmptyPassword(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "shed-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	configContent := `[shed-db]
password = ""
`

	mkdir(t, tmpDir)
	writeConfigFile(t, tmpDir, configContent)
	writeDatabaseFile(t, tmpDir)

	if validatePath(tmpDir) {
		t.Error("validatePath returned true when password is empty")
	}
}

// Helper function to create a valid test environment at the specified path.
func setupValidPath(t *testing.T, path string) {
	t.Helper()

	mkdir(t, path)

	// Create SQLite database file
	dbPath := filepath.Join(path, "shed.db")
	if err := os.WriteFile(dbPath, []byte{}, 0o644); err != nil {
		t.Fatalf("Failed to create db file: %v", err)
	}

	// Create valid config.toml
	configPath := filepath.Join(path, "config.toml")

	configContent := `[shed-db]
password = "test_password_123"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
}

func clearBuffer() {
	buf.Reset()
}

func createValidPathCWD(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	p := path.Join(dir, ".shed")

	createValidPath(t, p)

	return p
}

func createValidPath(t *testing.T, p string) {
	t.Helper()

	mkdir(t, p)

	// Create SQLite database file
	dbPath := filepath.Join(p, "shed.db")
	if err := os.WriteFile(dbPath, []byte{}, 0o644); err != nil {
		t.Fatalf("Failed to create db file: %v", err)
	}

	// Create valid config.toml
	configPath := filepath.Join(p, "config.toml")

	configContent := `[shed-db]
password = "test_password_123"
`
	if err := os.WriteFile(configPath, []byte(configContent), 0o644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}

	t.Cleanup(func() {
		// Clear the buffer after each test
		buf.Reset()
	})
}

func mkdir(t *testing.T, p string) {
	t.Helper()

	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatalf("Failed to create directory %q: %v", p, err)
	}

	t.Cleanup(func() {
		if err := os.RemoveAll(p); err != nil {
			t.Fatalf("Failed to clean up directory %q: %v", p, err)
		}
	})
}

func writeConfigFile(t *testing.T, dir, content string) {
	t.Helper()

	// Create only config.toml, not the database
	configPath := filepath.Join(dir, "config.toml")

	if err := os.WriteFile(configPath, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to create config file: %v", err)
	}
}

func writeDatabaseFile(t *testing.T, dir string) {
	t.Helper()

	dbPath := filepath.Join(dir, "shed.db")
	if err := os.WriteFile(dbPath, []byte{}, 0o644); err != nil {
		t.Fatalf("Failed to create db file: %v", err)
	}
}

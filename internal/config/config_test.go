package config

import (
	"bytes"
	"errors"
	"os"
	"path"
	"testing"

	"h3jfc/shed/internal/logger"
)

var buf bytes.Buffer

func init() {
	logger.SetWriter(&buf)
	logger.SetMode(logger.ModeVerbose)
	logger.New()
}

func TestFindPath_WithEnvironmentVariable(t *testing.T) {
	expectedPath := createValidPath(t)
	t.Setenv("SHED_DIR", expectedPath)

	// Execute
	result, err := FindPath()
	if err != nil {
		t.Fatalf("Expected FindPath() to succeed, but got error: %v", err)
	}

	// Assert
	if result != expectedPath {
		t.Errorf("Expected FindPath() to return %q, but got %q", expectedPath, result)
	}

	if !bytes.Contains(buf.Bytes(), []byte("Found SHED_DIR environment variable")) {
		t.Errorf("Expected debug log about SHED_DIR environment variable, but it was not found")
	}
}

func TestFindPath_NoEnv_DefaultPathFound(t *testing.T) {
	t.Cleanup(clearBuffer)
	t.Setenv("SHED_DIR", "")

	for _, exp := range DefaultConfigPaths { // nolint:paralleltest
		t.Run("test_dir_"+exp, func(t *testing.T) {
			createValidDefaultPath(t, exp)

			// Execute
			result, err := FindPath()
			if err != nil {
				t.Fatalf("Expected FindPath() to succeed, but got error: %v", err)
			}

			// Assert
			if result != exp {
				t.Errorf("Expected FindPath() to return %s, but got %q", exp, result)
			}

			if !bytes.Contains(buf.Bytes(), []byte("SHED_DIR environment variable not set")) {
				t.Errorf("Expected debug log about SHED_DIR environment variable not set, but it was not found")
			}
		})
	}
}

func TestFindPath_NothingFound(t *testing.T) {
	t.Cleanup(clearBuffer)
	t.Setenv("SHED_DIR", "")

	// Execute
	result, err := FindPath()
	if err == nil {
		t.Fatalf("Expected FindPath() to return error, but got none")
	}

	if !errors.Is(err, ErrNoPathFound) {
		t.Fatalf("Expected FindPath() to return ErrNoPathFound, but got: %v", err)
	}

	// Assert
	if result != "" {
		t.Errorf("Expected FindPath() to return empty string, but got %q", result)
	}

	if !bytes.Contains(buf.Bytes(), []byte("SHED_DIR environment variable not set")) {
		t.Errorf("Expected debug log about SHED_DIR environment variable not set, but it was not found")
	}

	if !bytes.Contains(buf.Bytes(), []byte("No existing shed path found.")) {
		t.Errorf("Expected debug log about no existing shed path found, but it was not found")
	}
}

func TestShedPath_WithEnvironmentVariable(t *testing.T) {
	expectedPath := createValidPath(t)
	t.Setenv("SHED_DIR", expectedPath)

	// Execute
	result := shedPath()

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
	result := shedPath()

	// Assert
	if result != "" {
		t.Errorf("Expected shedPath() to return empty string, but got %q", result)
	}

	if !bytes.Contains(buf.Bytes(), []byte("SHED_DIR environment variable not set")) {
		t.Errorf("Expected debug log about SHED_DIR environment variable not set, but it was not found")
	}
}

func clearBuffer() {
	buf.Reset()
}

func createValidPath(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}

	p := path.Join(dir, ".shed")

	t.Cleanup(func() {
		// Clear the buffer after each test
		buf.Reset()

		// Clean up the created test path
		if err := os.RemoveAll(p); err != nil {
			t.Fatalf("Failed to clean up test path %q: %v", p, err)
		}
	})

	return p
}

func createValidDefaultPath(t *testing.T, p string) {
	t.Helper()

	if err := os.MkdirAll(p, 0o755); err != nil {
		t.Fatalf("Failed to create test path %q: %v", p, err)
	}

	t.Cleanup(func() {
		// Clear the buffer after each test
		buf.Reset()

		// Clean up the created test path
		if err := os.RemoveAll(p); err != nil {
			t.Fatalf("Failed to clean up test path %q: %v", p, err)
		}
	})
}

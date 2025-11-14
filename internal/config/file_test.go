package config

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCreateConfigFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "test_password_123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // Function doesn't validate, just writes
		},
		{
			name:     "special characters in password",
			password: `p@ssw0rd!"#$%`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()

			err := createConfigFile(tmpDir, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("createConfigFile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				verifyConfigFile(t, tmpDir, tt.password)
			}
		})
	}
}

func verifyConfigFile(t *testing.T, tmpDir, password string) {
	t.Helper()

	configPath := filepath.Join(tmpDir, defaultConfigName)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file was not created")
	}

	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("failed to stat config file: %v", err)
	}

	// Skip permission check on Windows as it doesn't support Unix-style permissions
	if runtime.GOOS != "windows" {
		if info.Mode().Perm() != defaultFilePerms {
			t.Errorf("incorrect file permissions: got %v, want %v", info.Mode().Perm(), defaultFilePerms)
		}
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	if len(password) > 0 && !contains(string(content), password) {
		t.Error("config file does not contain the password")
	}
}

func TestCreateEmptyFile(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(string) string // Returns the file path to create
		wantErr bool
	}{
		{
			name: "create file in valid directory",
			setup: func(dir string) string {
				return filepath.Join(dir, "test.db")
			},
			wantErr: false,
		},
		{
			name: "create file in non-existent directory",
			setup: func(dir string) string {
				return filepath.Join(dir, "nonexistent", "test.db")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tmpDir := t.TempDir()
			filePath := tt.setup(tmpDir)

			err := createEmptyFile(filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("createEmptyFile() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !tt.wantErr {
				// Verify file exists
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Error("file was not created")
				}

				// Verify file is empty
				info, err := os.Stat(filePath)
				if err != nil {
					t.Fatalf("failed to stat file: %v", err)
				}

				if info.Size() != 0 {
					t.Errorf("file is not empty: size = %d", info.Size())
				}
			}
		})
	}
}

func TestValidatePath(t *testing.T) {
	t.Parallel()

	for _, tt := range getValidatePathTests() {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			path := tt.setup(t)

			got := validatePath(path)
			if got != tt.want {
				t.Errorf("validatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

// nolint:funlen // Test case table is long but necessary for thorough coverage
func getValidatePathTests() []struct {
	name  string
	setup func(*testing.T) string
	want  bool
} {
	return []struct {
		name  string
		setup func(*testing.T) string
		want  bool
	}{
		{
			name:  "valid shed directory",
			setup: setupValidShedDir,
			want:  true,
		},
		{
			name: "missing database file",
			setup: func(t *testing.T) string {
				t.Helper()

				return setupDirWithConfig(t, "test_password")
			},
			want: false,
		},
		{
			name:  "missing config file",
			setup: setupDirWithDB,
			want:  false,
		},
		{
			name: "config missing password field",
			setup: func(t *testing.T) string {
				t.Helper()

				return setupDirWithInvalidConfig(t, `[shed-db]
# missing password field

[settings]
`)
			},
			want: false,
		},
		{
			name: "config with empty password",
			setup: func(t *testing.T) string {
				t.Helper()

				tmpDir := t.TempDir()

				if err := createEmptyFile(filepath.Join(tmpDir, defaultDBName)); err != nil {
					t.Fatalf("failed to create db file: %v", err)
				}

				if err := createConfigFile(tmpDir, ""); err != nil {
					t.Fatalf("failed to create config file: %v", err)
				}

				return tmpDir
			},
			want: false,
		},
		{
			name: "path is a file not directory",
			setup: func(t *testing.T) string {
				t.Helper()

				tmpDir := t.TempDir()

				filePath := filepath.Join(tmpDir, "notadir.txt")
				if err := createEmptyFile(filePath); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}

				return filePath
			},
			want: false,
		},
		{
			name: "non-existent path",
			setup: func(t *testing.T) string {
				t.Helper()

				return "/nonexistent/path/that/does/not/exist"
			},
			want: false,
		},
		{
			name: "invalid toml syntax",
			setup: func(t *testing.T) string {
				t.Helper()

				return setupDirWithInvalidConfig(t, `[shed_db
password = "test"  # Missing closing bracket
`)
			},
			want: false,
		},
	}
}

func TestPromptUserForLocation(t *testing.T) {
	t.Parallel()

	// Note: This function requires stdin input, so we test error cases
	tests := []struct {
		name      string
		locations []string
		wantErr   bool
	}{
		{
			name:      "empty locations list",
			locations: []string{},
			wantErr:   true,
		},
		{
			name:      "nil locations list",
			locations: nil,
			wantErr:   true,
		},
		{
			name:      "valid locations list",
			locations: []string{"/home/user/.config", "/etc/shed"},
			wantErr:   false, // Would need user input to complete
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if len(tt.locations) == 0 {
				_, err := promptUserForLocation(tt.locations)
				if (err != nil) != tt.wantErr {
					t.Errorf("promptUserForLocation() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
			// For non-empty lists, we can't easily test without mocking stdin
		})
	}
}

// setupValidShedDir creates a valid shed directory with db and config.
func setupValidShedDir(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	dbPath := filepath.Join(tmpDir, defaultDBName)
	if err := createEmptyFile(dbPath); err != nil {
		t.Fatalf("failed to create db file: %v", err)
	}

	if err := createConfigFile(tmpDir, "test_password"); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	return tmpDir
}

// setupDirWithConfig creates a directory with only a config file.
func setupDirWithConfig(t *testing.T, password string) string {
	t.Helper()

	tmpDir := t.TempDir()

	if err := createConfigFile(tmpDir, password); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	return tmpDir
}

// setupDirWithDB creates a directory with only a database file.
func setupDirWithDB(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	dbPath := filepath.Join(tmpDir, defaultDBName)
	if err := createEmptyFile(dbPath); err != nil {
		t.Fatalf("failed to create db file: %v", err)
	}

	return tmpDir
}

// setupDirWithInvalidConfig creates a directory with db and invalid config.
func setupDirWithInvalidConfig(t *testing.T, content string) string {
	t.Helper()

	tmpDir := t.TempDir()

	dbPath := filepath.Join(tmpDir, defaultDBName)
	if err := createEmptyFile(dbPath); err != nil {
		t.Fatalf("failed to create db file: %v", err)
	}

	configPath := filepath.Join(tmpDir, defaultConfigName)
	if err := os.WriteFile(configPath, []byte(content), defaultFilePerms); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	return tmpDir
}

// Helper function for string contains check.
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

package os

import "testing"

func TestToString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    OS
		expected string
	}{
		{
			name:     "Darwin OS",
			input:    Darwin,
			expected: "darwin",
		},
		{
			name:     "Linux OS",
			input:    Linux,
			expected: "linux",
		},
		{
			name:     "Windows OS",
			input:    Windows,
			expected: "windows",
		},
		{
			name:     "Unknown OS",
			input:    Unknown,
			expected: "unknown",
		},
		{
			name:     "Invalid OS value defaults to unknown",
			input:    OS(999),
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ToString(tt.input)
			if result != tt.expected {
				t.Errorf("ToString(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToOS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		expected OS
	}{
		{
			name:     "darwin string to Darwin OS",
			input:    "darwin",
			expected: Darwin,
		},
		{
			name:     "linux string to Linux OS",
			input:    "linux",
			expected: Linux,
		},
		{
			name:     "windows string to Windows OS",
			input:    "windows",
			expected: Windows,
		},
		{
			name:     "unknown string defaults to Unknown OS",
			input:    "unknown",
			expected: Unknown,
		},
		{
			name:     "invalid string defaults to Unknown OS",
			input:    "freebsd",
			expected: Unknown,
		},
		{
			name:     "empty string defaults to Unknown OS",
			input:    "",
			expected: Unknown,
		},
		{
			name:     "case sensitive - Darwin capitalized defaults to Unknown",
			input:    "Darwin",
			expected: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := ToOS(tt.input)
			if result != tt.expected {
				t.Errorf("ToOS(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestOSConstants(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		os       OS
		expected int
	}{
		{
			name:     "Darwin is 0",
			os:       Darwin,
			expected: 0,
		},
		{
			name:     "Linux is 1",
			os:       Linux,
			expected: 1,
		},
		{
			name:     "Windows is 2",
			os:       Windows,
			expected: 2,
		},
		{
			name:     "Unknown is 3",
			os:       Unknown,
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if int(tt.os) != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, int(tt.os), tt.expected)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		os   OS
	}{
		{
			name: "Darwin round trip",
			os:   Darwin,
		},
		{
			name: "Linux round trip",
			os:   Linux,
		},
		{
			name: "Windows round trip",
			os:   Windows,
		},
		{
			name: "Unknown round trip",
			os:   Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			str := ToString(tt.os)

			result := ToOS(str)
			if result != tt.os {
				t.Errorf("Round trip failed: OS(%v) -> %q -> OS(%v)", tt.os, str, result)
			}
		})
	}
}

func TestStringRoundTrip(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		str  string
	}{
		{
			name: "darwin string round trip",
			str:  "darwin",
		},
		{
			name: "linux string round trip",
			str:  "linux",
		},
		{
			name: "windows string round trip",
			str:  "windows",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			os := ToOS(tt.str)

			result := ToString(os)
			if result != tt.str {
				t.Errorf("String round trip failed: %q -> OS(%v) -> %q", tt.str, os, result)
			}
		})
	}
}

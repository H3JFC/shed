package logfmt

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

var messageLevelTests = []struct {
	name        string
	level       slog.Level
	msg         string
	attrs       []slog.Attr
	contains    []string
	notContains []string
}{
	{
		name:  "info level",
		level: slog.LevelInfo,
		msg:   "test message",
		contains: []string{
			"INFO",
			"test message",
		},
		notContains: []string{
			"][", // No timestamp/level brackets in message-level mode
		},
	},
	{
		name:  "warn level",
		level: slog.LevelWarn,
		msg:   "warning message",
		contains: []string{
			"WARN",
			"warning message",
		},
	},
	{
		name:  "error level",
		level: slog.LevelError,
		msg:   "error occurred",
		contains: []string{
			"ERROR",
			"error occurred",
		},
	},
	{
		name:  "debug level",
		level: slog.LevelDebug,
		msg:   "debug info",
		contains: []string{
			"DEBUG",
			"debug info",
		},
	},
	{
		name:  "with attributes",
		level: slog.LevelWarn,
		msg:   "formatting file",
		attrs: []slog.Attr{
			slog.String("file", "main.go"),
			slog.Int("line", 42),
		},
		contains: []string{
			"WARN",
			"formatting file",
			"file=main.go",
			"line=42",
		},
	},
}

var messageOnlyTests = []struct {
	name        string
	level       slog.Level
	msg         string
	attrs       []slog.Attr
	contains    []string
	notContains []string
}{
	{
		name:  "message-only mode no level",
		level: slog.LevelInfo,
		msg:   "simple message",
		contains: []string{
			"simple message",
		},
		notContains: []string{
			"INFO",
			"][", // No timestamp/level brackets
			"T",  // No RFC3339 timestamp
		},
	},
	{
		name:  "message-only mode no timestamp",
		level: slog.LevelWarn,
		msg:   "warning text",
		contains: []string{
			"warning text",
		},
		notContains: []string{
			"WARN",
			"T",  // No RFC3339 timestamp
			"][", // No timestamp/level brackets
		},
	},
	{
		name:  "message-only mode with attributes",
		level: slog.LevelError,
		msg:   "error text",
		attrs: []slog.Attr{
			slog.String("code", "E001"),
		},
		contains: []string{
			"error text",
			"code=E001",
		},
		notContains: []string{
			"ERROR",
		},
	},
}

func TestCustomHandler_ModeMessageLevel(t *testing.T) {
	t.Parallel()

	for _, tt := range messageLevelTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			handler := NewCustomHandler(&buf, ModeMessageLevel)
			logger := slog.New(handler)

			logMessage(logger, tt.level, tt.msg, tt.attrs)

			output := buf.String()

			assertContains(t, output, tt.contains)
			assertNotContains(t, output, tt.notContains)
		})
	}
}

func TestCustomHandler_ModeVerbose(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		level    slog.Level
		msg      string
		attrs    []slog.Attr
		contains []string
	}{
		{
			name:  "verbose with timestamp",
			level: slog.LevelInfo,
			msg:   "verbose message",
			contains: []string{
				"[",
				"]",
				"INFO",
				"verbose message",
			},
		},
		{
			name:  "verbose with attributes",
			level: slog.LevelError,
			msg:   "error message",
			attrs: []slog.Attr{
				slog.String("error", "file not found"),
			},
			contains: []string{
				"ERROR",
				"error message",
				"error=file not found",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			handler := NewCustomHandler(&buf, ModeVerbose)
			logger := slog.New(handler)

			logMessage(logger, tt.level, tt.msg, tt.attrs)

			output := buf.String()

			assertContains(t, output, tt.contains)
			// Verify timestamp format (RFC3339Nano contains T and Z or timezone offset)
			if !strings.Contains(output, "T") {
				t.Errorf("Expected output to contain timestamp with 'T' separator (RFC3339 format), got: %s", output)
			}
		})
	}
}

func TestCustomHandler_ModeMessageOnly(t *testing.T) {
	t.Parallel()

	for _, tt := range messageOnlyTests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			handler := NewCustomHandler(&buf, ModeMessageOnly)
			logger := slog.New(handler)

			logMessage(logger, tt.level, tt.msg, tt.attrs)

			output := buf.String()

			assertContains(t, output, tt.contains)
			assertNotContains(t, output, tt.notContains)
		})
	}
}

func TestCustomHandler_ColorCodes(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	handler := NewCustomHandler(&buf, ModeMessageLevel)
	logger := slog.New(handler)

	tests := []struct {
		name      string
		logFunc   func(string, ...any)
		colorCode string
		levelName string
	}{
		{"warn has yellow", logger.Warn, colorYellow, "WARN"},
		{"error has red", logger.Error, colorRed, "ERROR"},
		{"info has blue", logger.Info, colorBlue, "INFO"},
		{"debug has cyan", logger.Debug, colorCyan, "DEBUG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf.Reset()
			tt.logFunc("test message")

			output := buf.String()

			if !strings.Contains(output, tt.colorCode) {
				t.Errorf("Expected output to contain color code for %s, got: %s", tt.levelName, output)
			}

			if !strings.Contains(output, colorReset) {
				t.Errorf("Expected output to contain color reset code, got: %s", output)
			}
		})
	}
}

func TestCustomHandler_Enabled(t *testing.T) {
	t.Parallel()

	handler := NewCustomHandler(&bytes.Buffer{}, ModeMessageLevel)
	ctx := context.Background()

	// All levels should be enabled
	levels := []slog.Level{
		slog.LevelDebug,
		slog.LevelInfo,
		slog.LevelWarn,
		slog.LevelError,
	}

	for _, level := range levels {
		if !handler.Enabled(ctx, level) {
			t.Errorf("Expected level %v to be enabled", level)
		}
	}
}

func TestCustomHandler_WithAttrs(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	handler := NewCustomHandler(&buf, ModeMessageLevel)

	// Add attributes to handler
	handlerWithAttrs := handler.WithAttrs([]slog.Attr{
		slog.String("component", "test"),
	})

	// Verify it returns a new handler
	if handlerWithAttrs == handler {
		t.Error("WithAttrs should return a new handler instance")
	}

	// Verify it's still a CustomHandler
	if _, ok := handlerWithAttrs.(*CustomHandler); !ok {
		t.Error("WithAttrs should return a *CustomHandler")
	}
}

func TestCustomHandler_WithGroup(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer

	handler := NewCustomHandler(&buf, ModeMessageLevel)

	// Add group to handler
	handlerWithGroup := handler.WithGroup("test-group")

	// Verify it returns a new handler
	if handlerWithGroup == handler {
		t.Error("WithGroup should return a new handler instance")
	}

	// Verify it's still a CustomHandler
	if _, ok := handlerWithGroup.(*CustomHandler); !ok {
		t.Error("WithGroup should return a *CustomHandler")
	}
}

func TestCustomHandler_AllModes(t *testing.T) {
	t.Parallel()

	modes := []LogMode{ModeMessageLevel, ModeVerbose, ModeMessageOnly}

	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer

			handler := NewCustomHandler(&buf, mode)
			logger := slog.New(handler)

			logger.Info("test message", "key", "value")

			output := buf.String()
			if output == "" {
				t.Errorf("Expected output for mode %s, got empty string", mode)
			}

			// All modes should contain the message
			if !strings.Contains(output, "test message") {
				t.Errorf("Expected output to contain 'test message' for mode %s, got: %s", mode, output)
			}
		})
	}
}

// Helper function to convert []slog.Attr to []any for logger methods.
func convertAttrs(attrs []slog.Attr) []any {
	result := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		result = append(result, attr.Key, attr.Value.Any())
	}

	return result
}

// logMessage logs a message at the specified level with optional attributes.
func logMessage(logger *slog.Logger, level slog.Level, msg string, attrs []slog.Attr) {
	switch level {
	case slog.LevelDebug:
		logger.Debug(msg, convertAttrs(attrs)...)
	case slog.LevelInfo:
		logger.Info(msg, convertAttrs(attrs)...)
	case slog.LevelWarn:
		logger.Warn(msg, convertAttrs(attrs)...)
	case slog.LevelError:
		logger.Error(msg, convertAttrs(attrs)...)
	}
}

// assertContains checks that the output contains all expected strings.
func assertContains(t *testing.T, output string, expected []string) {
	t.Helper()

	for _, exp := range expected {
		if !strings.Contains(output, exp) {
			t.Errorf("Expected output to contain %q, got: %s", exp, output)
		}
	}
}

// assertNotContains checks that the output does not contain any unexpected strings.
func assertNotContains(t *testing.T, output string, unexpected []string) {
	t.Helper()

	for _, unexp := range unexpected {
		if strings.Contains(output, unexp) {
			t.Errorf("Expected output NOT to contain %q, got: %s", unexp, output)
		}
	}
}

// Benchmark tests.
func BenchmarkCustomHandler_MessageLevel(b *testing.B) {
	var buf bytes.Buffer

	handler := NewCustomHandler(&buf, ModeMessageLevel)
	logger := slog.New(handler)

	b.ResetTimer()

	for i := range b.N {
		logger.Info("benchmark message", "iteration", i)
		buf.Reset()
	}
}

func BenchmarkCustomHandler_Verbose(b *testing.B) {
	var buf bytes.Buffer

	handler := NewCustomHandler(&buf, ModeVerbose)
	logger := slog.New(handler)

	b.ResetTimer()

	for i := range b.N {
		logger.Info("benchmark message", "iteration", i)
		buf.Reset()
	}
}

func BenchmarkCustomHandler_MessageOnly(b *testing.B) {
	var buf bytes.Buffer

	handler := NewCustomHandler(&buf, ModeMessageOnly)
	logger := slog.New(handler)

	b.ResetTimer()

	for i := range b.N {
		logger.Info("benchmark message", "iteration", i)
		buf.Reset()
	}
}

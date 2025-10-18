package logfmt

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
)

func TestCustomHandler_ModeMessageLevel(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := NewCustomHandler(&buf, ModeMessageLevel)
			logger := slog.New(handler)

			switch tt.level {
			case slog.LevelDebug:
				logger.Debug(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelInfo:
				logger.Info(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelWarn:
				logger.Warn(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelError:
				logger.Error(tt.msg, convertAttrs(tt.attrs)...)
			}

			output := buf.String()

			// Check for expected content
			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}

			// Check for unexpected content
			for _, unexpected := range tt.notContains {
				if strings.Contains(output, unexpected) {
					t.Errorf("Expected output NOT to contain %q, got: %s", unexpected, output)
				}
			}
		})
	}
}

func TestCustomHandler_ModeVerbose(t *testing.T) {
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
			var buf bytes.Buffer
			handler := NewCustomHandler(&buf, ModeVerbose)
			logger := slog.New(handler)

			switch tt.level {
			case slog.LevelDebug:
				logger.Debug(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelInfo:
				logger.Info(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelWarn:
				logger.Warn(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelError:
				logger.Error(tt.msg, convertAttrs(tt.attrs)...)
			}

			output := buf.String()

			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}

			// Verify timestamp format (RFC3339Nano contains T and Z or timezone offset)
			if !strings.Contains(output, "T") {
				t.Errorf("Expected output to contain timestamp with 'T' separator (RFC3339 format), got: %s", output)
			}
		})
	}
}

func TestCustomHandler_ModeMessageOnly(t *testing.T) {
	tests := []struct {
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := NewCustomHandler(&buf, ModeMessageOnly)
			logger := slog.New(handler)

			switch tt.level {
			case slog.LevelDebug:
				logger.Debug(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelInfo:
				logger.Info(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelWarn:
				logger.Warn(tt.msg, convertAttrs(tt.attrs)...)
			case slog.LevelError:
				logger.Error(tt.msg, convertAttrs(tt.attrs)...)
			}

			output := buf.String()

			for _, expected := range tt.contains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, got: %s", expected, output)
				}
			}

			for _, unexpected := range tt.notContains {
				if strings.Contains(output, unexpected) {
					t.Errorf("Expected output NOT to contain %q, got: %s", unexpected, output)
				}
			}
		})
	}
}

func TestCustomHandler_ColorCodes(t *testing.T) {
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
	modes := []LogMode{ModeMessageLevel, ModeVerbose, ModeMessageOnly}

	for _, mode := range modes {
		t.Run(string(mode), func(t *testing.T) {
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

// Helper function to convert []slog.Attr to []any for logger methods
func convertAttrs(attrs []slog.Attr) []any {
	result := make([]any, 0, len(attrs)*2)
	for _, attr := range attrs {
		result = append(result, attr.Key, attr.Value.Any())
	}
	return result
}

// Benchmark tests
func BenchmarkCustomHandler_MessageLevel(b *testing.B) {
	var buf bytes.Buffer
	handler := NewCustomHandler(&buf, ModeMessageLevel)
	logger := slog.New(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", "iteration", i)
		buf.Reset()
	}
}

func BenchmarkCustomHandler_Verbose(b *testing.B) {
	var buf bytes.Buffer
	handler := NewCustomHandler(&buf, ModeVerbose)
	logger := slog.New(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", "iteration", i)
		buf.Reset()
	}
}

func BenchmarkCustomHandler_MessageOnly(b *testing.B) {
	var buf bytes.Buffer
	handler := NewCustomHandler(&buf, ModeMessageOnly)
	logger := slog.New(handler)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", "iteration", i)
		buf.Reset()
	}
}

package logger

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

func TestNew_ReturnsSameInstance(t *testing.T) { // nolint:paralleltest
	Reset()

	logger1 := New(ModeMessageLevel)
	logger2 := New(ModeMessageLevel)

	if logger1 != logger2 {
		t.Error("New() should return the same instance")
	}
}

func TestNew_WritesToConfiguredWriter(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeMessageLevel)
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}
}

func TestModeVerbose_ShowsDebugMessages(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeVerbose)
	logger.Debug("debug message")
	logger.Info("info message")

	output := buf.String()
	if !strings.Contains(output, "debug message") {
		t.Errorf("Verbose mode should show debug messages, got: %s", output)
	}

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Verbose mode should show DEBUG level, got: %s", output)
	}

	if !strings.Contains(output, "info message") {
		t.Errorf("Verbose mode should show info messages, got: %s", output)
	}
}

func TestModeMessageLevel_FiltersDebugMessages(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeMessageLevel)
	logger.Debug("debug message")
	logger.Info("info message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Errorf("MessageLevel mode should NOT show debug messages, got: %s", output)
	}

	if !strings.Contains(output, "info message") {
		t.Errorf("MessageLevel mode should show info messages, got: %s", output)
	}

	if !strings.Contains(output, "INFO") {
		t.Errorf("MessageLevel mode should show INFO level, got: %s", output)
	}
}

func TestModeMessageOnly_FiltersDebugMessages(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeMessageOnly)
	logger.Debug("debug message")
	logger.Info("info message")

	output := buf.String()
	if strings.Contains(output, "debug message") {
		t.Errorf("MessageOnly mode should NOT show debug messages, got: %s", output)
	}

	if !strings.Contains(output, "info message") {
		t.Errorf("MessageOnly mode should show info messages, got: %s", output)
	}
}

func TestModeMessageLevel_ShowsLevelAndMessage(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeMessageLevel)
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "INFO") {
		t.Errorf("MessageLevel mode should include level, got: %s", output)
	}

	if !strings.Contains(output, "test message") {
		t.Errorf("MessageLevel mode should include message, got: %s", output)
	}

	// Should not contain timestamp in brackets
	if strings.Contains(output, "[") && strings.Contains(output, "]") {
		t.Errorf("MessageLevel mode should not contain timestamp in brackets, got: %s", output)
	}
}

func TestModeVerbose_ShowsTimestampLevelMessage(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeVerbose)
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "[") || !strings.Contains(output, "]") {
		t.Errorf("Verbose mode should include brackets for timestamp, got: %s", output)
	}

	if !strings.Contains(output, "INFO") {
		t.Errorf("Verbose mode should include level, got: %s", output)
	}

	if !strings.Contains(output, "test message") {
		t.Errorf("Verbose mode should include message, got: %s", output)
	}
}

func TestModeMessageOnly_ShowsOnlyMessage(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeMessageOnly)
	logger.Info("just message")

	output := buf.String()
	if !strings.Contains(output, "just message") {
		t.Errorf("MessageOnly mode should include message, got: %s", output)
	}

	// Should not contain INFO level or brackets
	if strings.Contains(output, "INFO") {
		t.Errorf("MessageOnly mode should not include level, got: %s", output)
	}

	if strings.Contains(output, "[") && strings.Contains(output, "]") {
		t.Errorf("MessageOnly mode should not include brackets, got: %s", output)
	}
}

func TestLoggerWithAttributes(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeMessageLevel)
	logger.Info("test message", "key", "value", "number", 42)

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Output should contain message, got: %s", output)
	}

	if !strings.Contains(output, "key=value") {
		t.Errorf("Output should contain key=value, got: %s", output)
	}

	if !strings.Contains(output, "number=42") {
		t.Errorf("Output should contain number=42, got: %s", output)
	}
}

func TestLogLevels_AllWork(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New(ModeVerbose)
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	output := buf.String()

	levels := []string{"DEBUG", "INFO", "WARN", "ERROR"}
	for _, level := range levels {
		if !strings.Contains(output, level) {
			t.Errorf("Output should contain %s level, got: %s", level, output)
		}
	}
}

func TestConcurrentNew(t *testing.T) { // nolint:paralleltest
	Reset()

	const goroutines = 100

	var wg sync.WaitGroup

	loggers := make([]*slog.Logger, goroutines)

	wg.Add(goroutines)

	for i := range goroutines {
		go func(index int) {
			defer wg.Done()

			loggers[index] = New(ModeMessageLevel)
		}(i)
	}

	wg.Wait()

	// All should be the same instance
	for i := range goroutines {
		if loggers[i] != loggers[0] {
			t.Error("Concurrent New() returned different instances")

			break
		}
	}
}

// Tests for package-level functions

func TestPackageLevel_Info(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	Info("package level info")

	output := buf.String()
	if !strings.Contains(output, "package level info") {
		t.Errorf("Expected output to contain 'package level info', got: %s", output)
	}

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected output to contain 'INFO', got: %s", output)
	}
}

func TestPackageLevel_Debug(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeVerbose)

	Debug("package level debug")

	output := buf.String()
	if !strings.Contains(output, "package level debug") {
		t.Errorf("Expected output to contain 'package level debug', got: %s", output)
	}

	if !strings.Contains(output, "DEBUG") {
		t.Errorf("Expected output to contain 'DEBUG', got: %s", output)
	}
}

func TestPackageLevel_Warn(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	Warn("package level warn")

	output := buf.String()
	if !strings.Contains(output, "package level warn") {
		t.Errorf("Expected output to contain 'package level warn', got: %s", output)
	}

	if !strings.Contains(output, "WARN") {
		t.Errorf("Expected output to contain 'WARN', got: %s", output)
	}
}

func TestPackageLevel_Error(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	Error("package level error")

	output := buf.String()
	if !strings.Contains(output, "package level error") {
		t.Errorf("Expected output to contain 'package level error', got: %s", output)
	}

	if !strings.Contains(output, "ERROR") {
		t.Errorf("Expected output to contain 'ERROR', got: %s", output)
	}
}

func TestPackageLevel_WithAttributes(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	Info("message with attrs", "user", "alice", "count", 10)

	output := buf.String()
	if !strings.Contains(output, "message with attrs") {
		t.Errorf("Expected output to contain message, got: %s", output)
	}

	if !strings.Contains(output, "user=alice") {
		t.Errorf("Expected output to contain user=alice, got: %s", output)
	}

	if !strings.Contains(output, "count=10") {
		t.Errorf("Expected output to contain count=10, got: %s", output)
	}
}

func TestPackageLevel_Context(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	ctx := context.Background()
	InfoContext(ctx, "context message")

	output := buf.String()
	if !strings.Contains(output, "context message") {
		t.Errorf("Expected output to contain 'context message', got: %s", output)
	}

	if !strings.Contains(output, "INFO") {
		t.Errorf("Expected output to contain 'INFO', got: %s", output)
	}
}

func TestPackageLevel_With(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	childLogger := With("component", "auth")
	childLogger.Info("authentication started")

	output := buf.String()
	if !strings.Contains(output, "authentication started") {
		t.Errorf("Expected output to contain message, got: %s", output)
	}

	if !strings.Contains(output, "component=auth") {
		t.Errorf("Expected output to contain component=auth, got: %s", output)
	}
}

func TestPackageLevel_WithGroup(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	groupLogger := WithGroup("request")
	groupLogger.Info("request received", "method", "GET")

	output := buf.String()
	if !strings.Contains(output, "request received") {
		t.Errorf("Expected output to contain message, got: %s", output)
	}
}

func TestPackageLevel_AlwaysUsesCurrentInstance(t *testing.T) { // nolint:paralleltest
	Reset()

	var buf1 bytes.Buffer
	SetWriter(&buf1)
	SetMode(ModeMessageLevel)

	// Initialize first instance
	Info("first message")

	// Create a new instance by resetting and reconfiguring
	Reset()

	var buf2 bytes.Buffer
	SetWriter(&buf2)
	SetMode(ModeMessageOnly)

	// Package-level function should use the new instance
	Info("second message")

	output1 := buf1.String()
	if !strings.Contains(output1, "first message") {
		t.Errorf("Expected first buffer to contain 'first message', got: %s", output1)
	}

	if !strings.Contains(output1, "INFO") {
		t.Errorf("Expected first buffer to contain 'INFO', got: %s", output1)
	}

	output2 := buf2.String()
	if !strings.Contains(output2, "second message") {
		t.Errorf("Expected second buffer to contain 'second message', got: %s", output2)
	}

	// Second should NOT contain INFO (message-only mode)
	if strings.Contains(output2, "INFO") {
		t.Errorf("Expected second buffer to NOT contain 'INFO' in message-only mode, got: %s", output2)
	}
}

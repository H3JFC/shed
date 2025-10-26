package logger

import (
	"bytes"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

func TestNew_ReturnsSameInstance(t *testing.T) {
	t.Parallel()
	Reset()

	logger1 := New()
	logger2 := New()

	if logger1 != logger2 {
		t.Error("New() should return the same instance")
	}
}

func TestNew_WritesToConfiguredWriter(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)

	logger := New()
	logger.Info("test message")

	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}
}

func TestModeVerbose_ShowsDebugMessages(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeVerbose)

	logger := New()
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

func TestModeMessageLevel_FiltersDebugMessages(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	logger := New()
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

func TestModeMessageOnly_FiltersDebugMessages(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageOnly)

	logger := New()
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

func TestModeMessageLevel_ShowsLevelAndMessage(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	logger := New()
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

func TestModeVerbose_ShowsTimestampLevelMessage(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeVerbose)

	logger := New()
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

func TestModeMessageOnly_ShowsOnlyMessage(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageOnly)

	logger := New()
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

func TestLoggerWithAttributes(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeMessageLevel)

	logger := New()
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

func TestLogLevels_AllWork(t *testing.T) {
	t.Parallel()
	Reset()

	var buf bytes.Buffer
	SetWriter(&buf)
	SetMode(ModeVerbose)

	logger := New()
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

			loggers[index] = New()
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

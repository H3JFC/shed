package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"
)

// LogMode defines the logging output format.
type (
	LogMode string
	Logger  = slog.Logger
)

const (
	ModeVerbose      LogMode = "verbose"
	ModeMessageLevel LogMode = "message-level"
	ModeMessageOnly  LogMode = "message-only"
)

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
)

var (
	instance *slog.Logger
	once     sync.Once
	mu       sync.RWMutex
	writer   io.Writer = os.Stdout
	mode               = ModeMessageLevel
)

// CustomHandler implements slog.Handler with different modes.
type CustomHandler struct {
	w     io.Writer
	mode  LogMode
	attrs []slog.Attr
	group string
}

// NewCustomHandler creates a new handler with the specified mode.
func NewCustomHandler(w io.Writer, mode LogMode) *CustomHandler {
	return &CustomHandler{
		w:    w,
		mode: mode,
	}
}

func (h *CustomHandler) Enabled(_ context.Context, level slog.Level) bool {
	// In verbose mode, enable debug and up
	if h.mode == ModeVerbose {
		return true
	}
	// In other modes, only Info and up
	return level >= slog.LevelInfo
}

func (h *CustomHandler) Handle(_ context.Context, r slog.Record) error {
	var output string

	switch h.mode {
	case ModeMessageLevel:
		output = h.formatMessageLevel(r)
	case ModeVerbose:
		output = h.formatVerbose(r)
	case ModeMessageOnly:
		output = h.formatMessageOnly(r)
	default:
		output = h.formatVerbose(r)
	}

	_, err := h.w.Write([]byte(output + "\n"))

	return err
}

func (h *CustomHandler) formatMessageLevel(r slog.Record) string {
	level := h.getLevelWithColor(r.Level)
	msg := r.Message
	// Add any additional attributes
	attrs := ""

	// Add handler-level attributes first
	for _, a := range h.attrs {
		if attrs != "" {
			attrs += " "
		}

		attrs += fmt.Sprintf("%s=%v", a.Key, a.Value)
	}

	r.Attrs(func(a slog.Attr) bool {
		if attrs != "" {
			attrs += " "
		}

		attrs += fmt.Sprintf("%s=%v", a.Key, a.Value)

		return true
	})

	if attrs != "" {
		return fmt.Sprintf("%s %s: %s", level, msg, attrs)
	}

	return fmt.Sprintf("%s %s", level, msg)
}

func (h *CustomHandler) formatVerbose(r slog.Record) string {
	timestamp := r.Time.Truncate(time.Microsecond).Format(time.RFC3339Nano)
	level := h.getLevelWithColor(r.Level)
	msg := r.Message
	// Add any additional attributes
	attrs := ""

	// Add handler-level attributes first
	for _, a := range h.attrs {
		if attrs != "" {
			attrs += " "
		}

		attrs += fmt.Sprintf("%s=%v", a.Key, a.Value)
	}

	r.Attrs(func(a slog.Attr) bool {
		if attrs != "" {
			attrs += " "
		}

		attrs += fmt.Sprintf("%s=%v", a.Key, a.Value)

		return true
	})

	if attrs != "" {
		return fmt.Sprintf("[%s][%s][%s] %s", timestamp, level, msg, attrs)
	}

	return fmt.Sprintf("[%s][%s][%s]", timestamp, level, msg)
}

func (h *CustomHandler) formatMessageOnly(r slog.Record) string {
	// Just the message, no level or timestamp
	msg := r.Message
	// Add any additional attributes
	attrs := ""

	// Add handler-level attributes first
	for _, a := range h.attrs {
		if attrs != "" {
			attrs += " "
		}

		attrs += fmt.Sprintf("%s=%v", a.Key, a.Value)
	}

	r.Attrs(func(a slog.Attr) bool {
		if attrs != "" {
			attrs += " "
		}

		attrs += fmt.Sprintf("%s=%v", a.Key, a.Value)

		return true
	})

	if attrs != "" {
		return fmt.Sprintf("%s %s", msg, attrs)
	}

	return msg
}

func (h *CustomHandler) getLevelWithColor(level slog.Level) string {
	switch level {
	case slog.LevelDebug:
		return colorCyan + "DEBUG" + colorReset
	case slog.LevelInfo:
		return colorBlue + "INFO" + colorReset
	case slog.LevelWarn:
		return colorYellow + "WARN" + colorReset
	case slog.LevelError:
		return colorRed + "ERROR" + colorReset
	default:
		return level.String()
	}
}

func (h *CustomHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(h.attrs)+len(attrs))
	copy(newAttrs, h.attrs)
	copy(newAttrs[len(h.attrs):], attrs)

	return &CustomHandler{
		w:     h.w,
		mode:  h.mode,
		attrs: newAttrs,
		group: h.group,
	}
}

func (h *CustomHandler) WithGroup(name string) slog.Handler {
	return &CustomHandler{
		w:     h.w,
		mode:  h.mode,
		attrs: h.attrs,
		group: name,
	}
}

func ModeFromString(s string) LogMode {
	switch s {
	case "verbose":
		return ModeVerbose
	case "message-level":
		return ModeMessageLevel
	case "message-only":
		return ModeMessageOnly
	default:
		return ModeMessageLevel
	}
}

// SetWriter sets the output writer (must be called before Get).
func SetWriter(w io.Writer) {
	mu.Lock()
	defer mu.Unlock()

	writer = w
}

// SetMode sets the logger mode (must be called before Get).
func SetMode(m LogMode) {
	mu.Lock()
	defer mu.Unlock()

	mode = m
}

func Get() *slog.Logger {
	mu.RLock()

	if instance != nil {
		mu.RUnlock()

		return instance
	}

	mu.RUnlock()

	// Need write lock to create instance
	mu.Lock()
	defer mu.Unlock()

	// Double-check after acquiring write lock
	if instance != nil {
		return instance
	}

	instance = slog.New(NewCustomHandler(writer, mode))

	return instance
}

// New returns the singleton logger instance.
func New() *slog.Logger {
	once.Do(func() {
		mu.RLock()

		w := writer
		m := mode

		mu.RUnlock()

		instance = slog.New(NewCustomHandler(w, m))
	})

	mu.RLock()
	defer mu.RUnlock()

	return instance
}

// Set allows setting a custom logger (useful for testing).
func Set(l *slog.Logger) {
	mu.Lock()
	defer mu.Unlock()

	instance = l
}

// Reset resets the singleton (useful for testing).
func Reset() {
	mu.Lock()
	defer mu.Unlock()

	instance = nil
	once = sync.Once{}
	writer = os.Stdout
	mode = ModeMessageLevel
}

// Package-level logging functions that always use the current instance.

// Debug logs a debug message with optional key-value pairs.
func Debug(msg string, args ...any) {
	Get().Debug(msg, args...)
}

// Info logs an info message with optional key-value pairs.
func Info(msg string, args ...any) {
	Get().Info(msg, args...)
}

// Warn logs a warning message with optional key-value pairs.
func Warn(msg string, args ...any) {
	Get().Warn(msg, args...)
}

// Error logs an error message with optional key-value pairs.
func Error(msg string, args ...any) {
	Get().Error(msg, args...)
}

// DebugContext logs a debug message with context and optional key-value pairs.
func DebugContext(ctx context.Context, msg string, args ...any) {
	Get().DebugContext(ctx, msg, args...)
}

// InfoContext logs an info message with context and optional key-value pairs.
func InfoContext(ctx context.Context, msg string, args ...any) {
	Get().InfoContext(ctx, msg, args...)
}

// WarnContext logs a warning message with context and optional key-value pairs.
func WarnContext(ctx context.Context, msg string, args ...any) {
	Get().WarnContext(ctx, msg, args...)
}

// ErrorContext logs an error message with context and optional key-value pairs.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	Get().ErrorContext(ctx, msg, args...)
}

// Log logs at the specified level with optional key-value pairs.
func Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	Get().Log(ctx, level, msg, args...)
}

// LogAttrs logs at the specified level with attributes.
func LogAttrs(ctx context.Context, level slog.Level, msg string, attrs ...slog.Attr) {
	Get().LogAttrs(ctx, level, msg, attrs...)
}

// With returns a Logger that includes the given attributes.
func With(args ...any) *slog.Logger {
	return Get().With(args...)
}

// WithGroup returns a Logger that starts a group.
func WithGroup(name string) *slog.Logger {
	return Get().WithGroup(name)
}

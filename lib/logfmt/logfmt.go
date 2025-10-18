package logfmt

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// LogMode defines the logging output format.
type LogMode string

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

func (h *CustomHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
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
	timestamp := r.Time.Format(time.RFC3339Nano)
	level := h.getLevelWithColor(r.Level)
	msg := r.Message

	// Add any additional attributes
	attrs := ""

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

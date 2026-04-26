package share

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"context"
	"fmt"
	"time"
)

type PrettyHandler struct {
	level slog.Level
}

func NewPrettyHandler(level slog.Level) *PrettyHandler {
	return &PrettyHandler{level: level}
}

// ANSI colors
const (
	reset  = "\033[0m"
	gray   = "\033[90m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
)

func levelColor(l slog.Level) string {
	switch {
	case l >= slog.LevelError:
		return red
	case l >= slog.LevelWarn:
		return yellow
	case l >= slog.LevelInfo:
		return green
	default:
		return gray
	}
}

func (h *PrettyHandler) Enabled(_ context.Context, l slog.Level) bool {
	return l >= h.level
}

func runtimeCallersFrame(pc uintptr) runtime.Frame {
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()
	return frame
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := time.Now().Format("2006-01-02 15:04:05")

	levelStr := r.Level.String()
	color := levelColor(r.Level)

	// source (file:line)
	source := ""
	if r.PC != 0 {
		fs := runtimeCallersFrame(r.PC)
		source = fmt.Sprintf("%s:%d", filepath.Base(fs.File), fs.Line)
	}

	// attributes
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		data, err := json.Marshal(a.Value.Any())
		if err == nil {
			attrs += fmt.Sprintf(" %s=%s", a.Key, string(data))
		} else {
			attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
		}
		return true
	})

	fmt.Fprintf(os.Stdout,
		color+"%s | %-5s | %s | %s | %s\n"+reset,
		timeStr,
		levelStr,
		source,
		r.Message,
		attrs,
	)

	return nil
}

func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyHandler) WithGroup(name string) slog.Handler {
	return h
}

// Opinioned set logging
func SetDefaultLog(level slog.Level) {
	handler := NewPrettyHandler(level)
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

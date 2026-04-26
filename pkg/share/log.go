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

const (
	ISO_8601 = "2026-04-26T05:28:56Z"
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

// LatencyColor is the ANSI color for latency
func latencyColor(latency time.Duration, defaultColor string) string {
	switch {
	case latency < 200 * time.Millisecond:
		return defaultColor
	case latency < time.Second:
		return yellow
	default:
		return red
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

// No caller data is show, thus limited the log data leak
// Only care about the call tree
func fullTracebackFrames(skip int) []runtime.Frame {
	// capture PCs; +10 is an upper bound for depth growth
	pcs := make([]uintptr, 64)
	n := runtime.Callers(2+skip, pcs) // skip runtime.Callers + this wrapper
	pcs = pcs[:n]

	frames := runtime.CallersFrames(pcs)
	var res []runtime.Frame
	for {
		frame, more := frames.Next()
		res = append(res, frame)
		if !more {
			break
		}
	}
	return res
}

func (h *PrettyHandler) Handle(_ context.Context, r slog.Record) error {
	// ISO 8601
	timeStr := time.Now().Format(ISO_8601)

	levelStr := r.Level.String()
	color := levelColor(r.Level)

	// source (file:line)
	source := ""
	if r.PC != 0 {
		fs := runtimeCallersFrame(r.PC)
		source = fmt.Sprintf("%s (%s:%d)", fs.Function, filepath.Base(fs.File), fs.Line)
	}

	var traceback []runtime.Frame = nil
	if r.Level >= slog.LevelError {
		traceback = fullTracebackFrames(0)
	}

	// attributes
	attrs := ""
	r.Attrs(func(a slog.Attr) bool {
		switch a.Value.Kind() {
		case slog.KindDuration:
			attrs += fmt.Sprintf(reset+latencyColor(a.Value.Duration(), color)+" %s=%8v"+reset+color, a.Key, a.Value.Duration())
		case slog.KindTime:
			attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Time().Format(ISO_8601))

		case slog.KindAny:
			fallthrough
		default:
			data, err := json.Marshal(a.Value.Any())
			if err != nil {
				attrs += fmt.Sprintf(" %s=%v", a.Key, a.Value.Any())
			} else {
				attrs += fmt.Sprintf(" %s=%s", a.Key, string(data))
			}

		}
		return true
	})

	fmt.Fprintf(os.Stdout,
		color+"%s | %-5s | %s | source=%s %s\n"+reset,
		timeStr,
		levelStr,
		r.Message,
		source,
		attrs,
	)

	for _, fs := range traceback {
		fmt.Printf("\t%s:%d\n", fs.File, fs.Line)
	}

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

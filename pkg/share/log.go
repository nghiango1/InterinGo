package share

import (
	"log/slog"
	"os"
	"path/filepath"
)

// Opinioned set logging
func SetDefaultLog(level slog.Level) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true, // include source file:line
		Level:     level,
		ReplaceAttr: func(gs []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				s := a.Value.Any().(*slog.Source)
				s.File = filepath.Base(s.File)
				return slog.Any(a.Key, s)
			}
			return a
		},
	}))
	slog.SetDefault(logger)
}

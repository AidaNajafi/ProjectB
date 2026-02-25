package logger

import (
	"log/slog"
	"os"
	"strings"
)

type Logger interface {
	Info(msg string, fields map[string]any)
	Error(msg string, fields map[string]any)
}

type SlogLogger struct {
	l *slog.Logger
}

func NewSlog(l *slog.Logger) *SlogLogger {
	return &SlogLogger{l: l}
}

func (s *SlogLogger) Info(msg string, fields map[string]any) {
	s.l.Info(msg, toArgs(fields)...)
}

func (s *SlogLogger) Error(msg string, fields map[string]any) {
	s.l.Error(msg, toArgs(fields)...)
}

func toArgs(fields map[string]any) []any {
	if len(fields) == 0 {
		return nil
	}
	args := make([]any, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return args
}

func ParseLevel(level string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func NewJSON(levelStr string) *SlogLogger {
	level := ParseLevel(levelStr)

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	return NewSlog(slog.New(handler))
}

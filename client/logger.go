package client

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"
)

// Logger is an interface for logging within the Admiral client.
// The SDK is silent by default (NoOpLogger). Use NewStdLogger or
// NewSlogLogger to enable output.
type Logger interface {
	// Debugf logs a debug message with printf-style formatting.
	Debugf(format string, args ...any)
	// Infof logs an info message with printf-style formatting.
	Infof(format string, args ...any)
	// Warnf logs a warning message with printf-style formatting.
	Warnf(format string, args ...any)
	// Errorf logs an error message with printf-style formatting.
	Errorf(format string, args ...any)
}

// Level represents a log severity level.
type Level int

const (
	// LevelDebug is the most verbose level.
	LevelDebug Level = iota
	// LevelInfo is the default level for StdLogger.
	LevelInfo
	// LevelWarn only logs warnings and errors.
	LevelWarn
	// LevelError only logs errors.
	LevelError
)

// String returns the human-readable name of the level.
func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	default:
		return fmt.Sprintf("LEVEL(%d)", int(l))
	}
}

// NoOpLogger is a logger that discards all log messages.
type NoOpLogger struct{}

// NewNoOpLogger creates a new NoOpLogger.
func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}

func (n *NoOpLogger) Debugf(format string, args ...any) {}
func (n *NoOpLogger) Infof(format string, args ...any)  {}
func (n *NoOpLogger) Warnf(format string, args ...any)  {}
func (n *NoOpLogger) Errorf(format string, args ...any) {}

// StdLogger writes log messages to an io.Writer at a configurable level.
type StdLogger struct {
	mu    sync.Mutex
	w     io.Writer
	level Level
}

// NewStdLogger creates a logger that writes to w at the given minimum level.
//
//	logger := client.NewStdLogger(os.Stderr, client.LevelInfo)
func NewStdLogger(w io.Writer, level Level) Logger {
	return &StdLogger{w: w, level: level}
}

func (s *StdLogger) Debugf(format string, args ...any) { s.logf(LevelDebug, format, args...) }
func (s *StdLogger) Infof(format string, args ...any)  { s.logf(LevelInfo, format, args...) }
func (s *StdLogger) Warnf(format string, args ...any)  { s.logf(LevelWarn, format, args...) }
func (s *StdLogger) Errorf(format string, args ...any) { s.logf(LevelError, format, args...) }

func (s *StdLogger) logf(level Level, format string, args ...any) {
	if level < s.level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	ts := time.Now().Format(time.RFC3339)

	s.mu.Lock()
	defer s.mu.Unlock()
	_, _ = fmt.Fprintf(s.w, "%s [%s] %s\n", ts, level, msg)
}

// SlogAdapter adapts a *slog.Logger to the Logger interface.
// Uses only the standard library — no external dependencies.
//
//	logger := client.NewSlogLogger(slog.Default())
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogLogger creates a Logger backed by a *slog.Logger.
func NewSlogLogger(logger *slog.Logger) Logger {
	return &SlogAdapter{logger: logger}
}

func (s *SlogAdapter) Debugf(format string, args ...any) {
	s.logger.Debug(fmt.Sprintf(format, args...))
}

func (s *SlogAdapter) Infof(format string, args ...any) {
	s.logger.Info(fmt.Sprintf(format, args...))
}

func (s *SlogAdapter) Warnf(format string, args ...any) {
	s.logger.Warn(fmt.Sprintf(format, args...))
}

func (s *SlogAdapter) Errorf(format string, args ...any) {
	s.logger.Error(fmt.Sprintf(format, args...))
}

package client

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/go-logr/logr"
)

// Logger is an interface for logging within the Admiral client.
// It supports structured logging with levels.
type Logger interface {
	// Debug logs a debug message with optional key-value pairs
	Debug(msg string, keysAndValues ...interface{})
	// Info logs an info message with optional key-value pairs
	Info(msg string, keysAndValues ...interface{})
	// Warn logs a warning message with optional key-value pairs
	Warn(msg string, keysAndValues ...interface{})
	// Error logs an error message with optional key-value pairs
	Error(err error, msg string, keysAndValues ...interface{})
	// WithValues returns a new Logger with additional key-value pairs
	WithValues(keysAndValues ...interface{}) Logger
	// WithName returns a new Logger with the specified name appended
	WithName(name string) Logger
}

// LogrAdapter adapts a logr.Logger to the Logger interface
type LogrAdapter struct {
	logger logr.Logger
}

// NewLogrAdapter creates a new LogrAdapter from a logr.Logger
func NewLogrAdapter(logger logr.Logger) Logger {
	return &LogrAdapter{logger: logger}
}

func (l *LogrAdapter) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.V(1).Info(msg, keysAndValues...)
}

func (l *LogrAdapter) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *LogrAdapter) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Info(msg, append([]interface{}{"level", "warning"}, keysAndValues...)...)
}

func (l *LogrAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	l.logger.Error(err, msg, keysAndValues...)
}

func (l *LogrAdapter) WithValues(keysAndValues ...interface{}) Logger {
	return &LogrAdapter{logger: l.logger.WithValues(keysAndValues...)}
}

func (l *LogrAdapter) WithName(name string) Logger {
	return &LogrAdapter{logger: l.logger.WithName(name)}
}

// SlogAdapter adapts a *slog.Logger to the Logger interface
type SlogAdapter struct {
	logger *slog.Logger
}

// NewSlogAdapter creates a new SlogAdapter from a *slog.Logger
func NewSlogAdapter(logger *slog.Logger) Logger {
	return &SlogAdapter{logger: logger}
}

func (s *SlogAdapter) Debug(msg string, keysAndValues ...interface{}) {
	s.logger.Debug(msg, keysAndValues...)
}

func (s *SlogAdapter) Info(msg string, keysAndValues ...interface{}) {
	s.logger.Info(msg, keysAndValues...)
}

func (s *SlogAdapter) Warn(msg string, keysAndValues ...interface{}) {
	s.logger.Warn(msg, keysAndValues...)
}

func (s *SlogAdapter) Error(err error, msg string, keysAndValues ...interface{}) {
	if err != nil {
		s.logger.Error(msg, append([]interface{}{"error", err}, keysAndValues...)...)
	} else {
		s.logger.Error(msg, keysAndValues...)
	}
}

func (s *SlogAdapter) WithValues(keysAndValues ...interface{}) Logger {
	// Convert key-value pairs to slog attributes and then to any slice
	var args []any
	for i := 0; i < len(keysAndValues); i += 2 {
		if i+1 < len(keysAndValues) {
			args = append(args, keysAndValues[i], keysAndValues[i+1])
		}
	}
	return &SlogAdapter{logger: s.logger.With(args...)}
}

func (s *SlogAdapter) WithName(name string) Logger {
	return &SlogAdapter{logger: s.logger.With("component", name)}
}

// DefaultLogger is a simple logger that uses the standard library log package
type DefaultLogger struct {
	prefix string
}

// NewDefaultLogger creates a new DefaultLogger
func NewDefaultLogger() Logger {
	return &DefaultLogger{}
}

func (d *DefaultLogger) Debug(msg string, keysAndValues ...interface{}) {
	// Debug is typically disabled in default logger
	if debugEnabled() {
		d.logf("DEBUG", msg, keysAndValues...)
	}
}

func (d *DefaultLogger) Info(msg string, keysAndValues ...interface{}) {
	d.logf("INFO", msg, keysAndValues...)
}

func (d *DefaultLogger) Warn(msg string, keysAndValues ...interface{}) {
	d.logf("WARN", msg, keysAndValues...)
}

func (d *DefaultLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if err != nil {
		d.logf("ERROR", fmt.Sprintf("%s: %v", msg, err), keysAndValues...)
	} else {
		d.logf("ERROR", msg, keysAndValues...)
	}
}

func (d *DefaultLogger) WithValues(keysAndValues ...interface{}) Logger {
	// For simplicity, default logger doesn't maintain context
	return d
}

func (d *DefaultLogger) WithName(name string) Logger {
	return &DefaultLogger{prefix: name}
}

func (d *DefaultLogger) logf(level, msg string, keysAndValues ...interface{}) {
	prefix := ""
	if d.prefix != "" {
		prefix = fmt.Sprintf("[%s] ", d.prefix)
	}

	if len(keysAndValues) > 0 {
		// Format key-value pairs
		var kvPairs []string
		for i := 0; i < len(keysAndValues); i += 2 {
			if i+1 < len(keysAndValues) {
				kvPairs = append(kvPairs, fmt.Sprintf("%v=%v", keysAndValues[i], keysAndValues[i+1]))
			}
		}
		log.Printf("%s[%s] %s %v", prefix, level, msg, kvPairs)
	} else {
		log.Printf("%s[%s] %s", prefix, level, msg)
	}
}

// debugEnabled checks if debug logging should be enabled
// This could be controlled by an environment variable
func debugEnabled() bool {
	// For now, disable debug in default logger
	// Could check DEBUG env var or similar
	return false
}

// NoOpLogger is a logger that discards all log messages
type NoOpLogger struct{}

// NewNoOpLogger creates a new NoOpLogger
func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}

func (n *NoOpLogger) Debug(msg string, keysAndValues ...interface{})            {}
func (n *NoOpLogger) Info(msg string, keysAndValues ...interface{})             {}
func (n *NoOpLogger) Warn(msg string, keysAndValues ...interface{})             {}
func (n *NoOpLogger) Error(err error, msg string, keysAndValues ...interface{}) {}
func (n *NoOpLogger) WithValues(keysAndValues ...interface{}) Logger            { return n }
func (n *NoOpLogger) WithName(name string) Logger                               { return n }

package client

import (
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/go-logr/logr"
	"github.com/stretchr/testify/assert"
)

func TestLogrAdapter(t *testing.T) {
	// Create a test logr.Logger (using logr's test implementation)
	logger := logr.Discard()
	adapter := NewLogrAdapter(logger)

	// These should not panic
	adapter.Debug("debug message", "key", "value")
	adapter.Info("info message", "key", "value")
	adapter.Warn("warning message", "key", "value")
	adapter.Error(errors.New("test error"), "error message", "key", "value")

	// Test WithValues and WithName
	childLogger := adapter.WithValues("component", "test").WithName("admiral")
	assert.NotNil(t, childLogger)

	// Should not panic
	childLogger.Info("child logger message")
}

func TestSlogAdapter(t *testing.T) {
	// Create a test slog.Logger (discard handler)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	adapter := NewSlogAdapter(logger)

	// These should not panic
	adapter.Debug("debug message", "key", "value")
	adapter.Info("info message", "key", "value")
	adapter.Warn("warning message", "key", "value")
	adapter.Error(errors.New("test error"), "error message", "key", "value")
	adapter.Error(nil, "error without err", "key", "value")

	// Test WithValues and WithName
	childLogger := adapter.WithValues("component", "test").WithName("admiral")
	assert.NotNil(t, childLogger)

	// Should not panic
	childLogger.Info("child logger message")
}

func TestDefaultLogger(t *testing.T) {
	logger := NewDefaultLogger()

	// These should not panic
	logger.Debug("debug message", "key", "value")
	logger.Info("info message", "key", "value")
	logger.Warn("warning message", "key", "value")
	logger.Error(errors.New("test error"), "error message", "key", "value")
	logger.Error(nil, "error without err", "key", "value")

	// Test WithValues and WithName
	childLogger := logger.WithValues("component", "test").WithName("admiral")
	assert.NotNil(t, childLogger)

	// Should not panic
	childLogger.Info("child logger message")
}

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()

	// These should not panic and should do nothing
	logger.Debug("debug message", "key", "value")
	logger.Info("info message", "key", "value")
	logger.Warn("warning message", "key", "value")
	logger.Error(errors.New("test error"), "error message", "key", "value")

	// Test WithValues and WithName
	childLogger := logger.WithValues("component", "test").WithName("admiral")
	assert.NotNil(t, childLogger)

	// Should not panic
	childLogger.Info("child logger message")
}

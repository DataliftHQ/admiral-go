package client

import (
	"bytes"
	"io"
	"log/slog"
	"strings"
	"sync"
	"testing"
)

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()

	// Should not panic
	logger.Debugf("debug %s", "msg")
	logger.Infof("info %s", "msg")
	logger.Warnf("warn %s", "msg")
	logger.Errorf("error %s", "msg")
}

func TestStdLogger_LevelFiltering(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		wantHas  []string
		wantMiss []string
	}{
		{
			name:     "debug level logs everything",
			level:    LevelDebug,
			wantHas:  []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]"},
			wantMiss: nil,
		},
		{
			name:     "info level skips debug",
			level:    LevelInfo,
			wantHas:  []string{"[INFO]", "[WARN]", "[ERROR]"},
			wantMiss: []string{"[DEBUG]"},
		},
		{
			name:     "warn level skips debug and info",
			level:    LevelWarn,
			wantHas:  []string{"[WARN]", "[ERROR]"},
			wantMiss: []string{"[DEBUG]", "[INFO]"},
		},
		{
			name:     "error level only logs errors",
			level:    LevelError,
			wantHas:  []string{"[ERROR]"},
			wantMiss: []string{"[DEBUG]", "[INFO]", "[WARN]"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := NewStdLogger(&buf, tt.level)

			logger.Debugf("debug message")
			logger.Infof("info message")
			logger.Warnf("warn message")
			logger.Errorf("error message")

			output := buf.String()
			for _, want := range tt.wantHas {
				if !strings.Contains(output, want) {
					t.Errorf("output should contain %q, got:\n%s", want, output)
				}
			}
			for _, miss := range tt.wantMiss {
				if strings.Contains(output, miss) {
					t.Errorf("output should NOT contain %q, got:\n%s", miss, output)
				}
			}
		})
	}
}

func TestStdLogger_FormatsMessage(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStdLogger(&buf, LevelDebug)

	logger.Infof("host=%s port=%d", "localhost", 9443)

	output := buf.String()
	if !strings.Contains(output, "host=localhost port=9443") {
		t.Errorf("expected formatted message, got: %s", output)
	}
	if !strings.Contains(output, "[INFO]") {
		t.Errorf("expected [INFO] prefix, got: %s", output)
	}
}

func TestStdLogger_ConcurrentWrites(t *testing.T) {
	var buf bytes.Buffer
	logger := NewStdLogger(&buf, LevelDebug)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			logger.Infof("goroutine %d", n)
		}(i)
	}
	wg.Wait()

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 100 {
		t.Errorf("expected 100 lines, got %d", len(lines))
	}
}

func TestSlogAdapter(t *testing.T) {
	logger := NewSlogLogger(slog.New(slog.NewTextHandler(io.Discard, nil)))

	// Should not panic
	logger.Debugf("debug %s", "msg")
	logger.Infof("info %s", "msg")
	logger.Warnf("warn %s", "msg")
	logger.Errorf("error %s", "msg")
}

func TestSlogAdapter_DelegatesToSlog(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := NewSlogLogger(slog.New(handler))

	logger.Infof("hello %s", "world")

	output := buf.String()
	if !strings.Contains(output, "hello world") {
		t.Errorf("expected slog output to contain %q, got: %s", "hello world", output)
	}
	if !strings.Contains(output, "level=INFO") {
		t.Errorf("expected slog output to contain level=INFO, got: %s", output)
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{Level(99), "LEVEL(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %q, want %q", got, tt.want)
			}
		})
	}
}

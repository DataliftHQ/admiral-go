package client

import (
	"runtime"
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()

	// Version should be non-empty
	if v == "" {
		t.Error("Version() returned empty string")
	}

	// In development/test context, expect "dev"
	// When used as a dependency, would return actual version like "v1.2.3"
	t.Logf("Version: %s", v)
}

func TestClientUserAgent(t *testing.T) {
	ua := ClientUserAgent()

	// Should contain the module path
	if !strings.Contains(ua, "go.admiral.io/admiral-go") {
		t.Errorf("ClientUserAgent() missing module path, got: %s", ua)
	}

	// Should contain version
	if !strings.Contains(ua, Version()) {
		t.Errorf("ClientUserAgent() missing version, got: %s", ua)
	}

	// Should contain platform info
	if !strings.Contains(ua, runtime.GOOS) {
		t.Errorf("ClientUserAgent() missing GOOS, got: %s", ua)
	}

	if !strings.Contains(ua, runtime.GOARCH) {
		t.Errorf("ClientUserAgent() missing GOARCH, got: %s", ua)
	}

	// Should contain Go version
	if !strings.Contains(ua, runtime.Version()) {
		t.Errorf("ClientUserAgent() missing Go version, got: %s", ua)
	}

	t.Logf("User-Agent: %s", ua)
}

func TestModulePath(t *testing.T) {
	// Verify the module path constant is correctly templated
	expected := "go.admiral.io/admiral-go"
	if modulePath != expected {
		t.Errorf("modulePath = %q, want %q", modulePath, expected)
	}
}

func TestIsRealVersion(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		// Real versions
		{"v1.0.0", true},
		{"v1.2.3", true},
		{"v0.1.0", true},
		{"v2.0.0-beta.1", true},

		// Development/pseudo-versions
		{"", false},
		{"(devel)", false},
		{"v0.0.0-20260106102548-abc123def456", false},
		{"v0.0.0-00010101000000-000000000000", false},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got := isRealVersion(tt.version)
			if got != tt.want {
				t.Errorf("isRealVersion(%q) = %v, want %v", tt.version, got, tt.want)
			}
		})
	}
}

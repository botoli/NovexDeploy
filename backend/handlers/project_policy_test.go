package handlers

import "testing"

func TestNormalizeProjectType(t *testing.T) {
	cases := map[string]string{
		"":             ProjectTypeBackend,
		"service":      ProjectTypeBackend,
		"backend":      ProjectTypeBackend,
		"telegram":     ProjectTypeTelegram,
		"telegram_bot": ProjectTypeTelegram,
	}
	for input, expected := range cases {
		if got := normalizeProjectType(input); got != expected {
			t.Fatalf("normalizeProjectType(%q) = %q, want %q", input, got, expected)
		}
	}
}

func TestIsFrontendDeployConfig(t *testing.T) {
	if !isFrontendDeployConfig("npm run build && vite build", "", "dist") {
		t.Fatal("expected frontend deploy config to be rejected")
	}
	if isFrontendDeployConfig("go build ./cmd/api", "./api", ".") {
		t.Fatal("expected backend deploy config to be allowed")
	}
}


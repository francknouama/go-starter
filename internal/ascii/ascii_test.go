package ascii

import (
	"strings"
	"testing"
)

func TestBanner(t *testing.T) {
	banner := Banner()

	// Test that banner is not empty
	if banner == "" {
		t.Error("Expected banner to not be empty")
	}

	// Test that banner contains ASCII art elements (checking for stylized patterns)
	if !strings.Contains(banner, "██") && !strings.Contains(banner, "╗") && !strings.Contains(banner, "║") {
		t.Error("Expected banner to contain ASCII art characters (box drawing)")
	}

	// Test that banner has multiple lines (ASCII art should be multi-line)
	lines := strings.Split(banner, "\n")
	if len(lines) < 3 {
		t.Errorf("Expected banner to have at least 3 lines, got %d", len(lines))
	}

	// Test that banner doesn't contain any control characters that might break terminals
	for i, line := range lines {
		for j, char := range line {
			// Allow printable ASCII, spaces, and common terminal-safe characters
			if char < 32 && char != 9 && char != 10 && char != 13 { // Allow tab, newline, carriage return
				t.Errorf("Banner contains non-printable character at line %d, position %d: %v", i, j, char)
			}
		}
	}
}

func TestBannerConsistency(t *testing.T) {
	// Test that Banner() returns the same result when called multiple times
	banner1 := Banner()
	banner2 := Banner()

	if banner1 != banner2 {
		t.Error("Expected Banner() to return consistent results")
	}
}

func TestBannerLength(t *testing.T) {
	banner := Banner()

	// Test that banner is reasonable length (not too short, not excessively long)
	if len(banner) < 50 {
		t.Error("Banner seems too short to be meaningful ASCII art")
	}

	if len(banner) > 5000 {
		t.Error("Banner seems excessively long")
	}
}

func TestBannerFormat(t *testing.T) {
	banner := Banner()
	lines := strings.Split(banner, "\n")

	// Test that no line is excessively long (ASCII art can be wide, so allow up to 200 chars)
	for i, line := range lines {
		if len(line) > 200 {
			t.Errorf("Line %d is too long (%d characters)", i, len(line))
		}
	}
}

// Test that the banner function is accessible and can be called without panicking
func TestBannerNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Banner() panicked: %v", r)
		}
	}()

	_ = Banner()
}

// Benchmark the Banner function to ensure it's performant
func BenchmarkBanner(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Banner()
	}
}

package ascii

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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
	for b.Loop() {
		Banner()
	}
}

// New comprehensive tests for standardized banner system

func TestBannerWithConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *BannerConfig
		expected string
	}{
		{
			name: "disabled banner",
			config: &BannerConfig{
				Enabled: false,
				Style:   StyleFull,
				Colors:  true,
			},
			expected: "",
		},
		{
			name: "none style banner",
			config: &BannerConfig{
				Enabled: true,
				Style:   StyleNone,
				Colors:  true,
			},
			expected: "",
		},
		{
			name: "minimal style banner",
			config: &BannerConfig{
				Enabled: true,
				Style:   StyleMinimal,
				Colors:  true,
			},
			expected: "STARTER", // Should contain part of logo
		},
		{
			name: "full style banner",
			config: &BannerConfig{
				Enabled: true,
				Style:   StyleFull,
				Colors:  true,
			},
			expected: "██", // Should contain ASCII art characters
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BannerWithConfig(tt.config)
			
			if tt.expected == "" {
				assert.Empty(t, result)
			} else {
				assert.Contains(t, result, tt.expected)
			}
		})
	}
}

func TestLogo(t *testing.T) {
	logo := Logo()
	
	assert.NotEmpty(t, logo)
	assert.Contains(t, logo, "STARTER")
}

func TestLogoWithConfig(t *testing.T) {
	tests := []struct {
		name     string
		config   *BannerConfig
		expected bool // whether logo should be present
	}{
		{
			name: "enabled logo",
			config: &BannerConfig{
				Enabled: true,
				Style:   StyleFull,
				Colors:  true,
			},
			expected: true,
		},
		{
			name: "disabled logo",
			config: &BannerConfig{
				Enabled: false,
				Style:   StyleFull,
				Colors:  true,
			},
			expected: false,
		},
		{
			name: "none style logo",
			config: &BannerConfig{
				Enabled: true,
				Style:   StyleNone,
				Colors:  true,
			},
			expected: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LogoWithConfig(tt.config)
			
			if tt.expected {
				assert.NotEmpty(t, result)
				assert.Contains(t, result, "STARTER")
			} else {
				assert.Empty(t, result)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	assert.True(t, config.Enabled)
	assert.Equal(t, StyleFull, config.Style)
	assert.True(t, config.Colors)
	assert.True(t, config.ShowOnHelp)
	assert.False(t, config.Quiet)
}

func TestConfigFromEnv(t *testing.T) {
	// Save original environment
	originalBanner := os.Getenv("GO_STARTER_BANNER")
	originalStyle := os.Getenv("GO_STARTER_BANNER_STYLE")
	originalNoColor := os.Getenv("NO_COLOR")
	
	// Cleanup after test
	defer func() {
		_ = os.Setenv("GO_STARTER_BANNER", originalBanner)
		_ = os.Setenv("GO_STARTER_BANNER_STYLE", originalStyle)
		_ = os.Setenv("NO_COLOR", originalNoColor)
	}()
	
	tests := []struct {
		name           string
		banner         string
		style          string
		noColor        string
		expectedEnabled bool
		expectedStyle   BannerStyle
		expectedColors  bool
	}{
		{
			name:           "default environment",
			banner:         "",
			style:          "",
			noColor:        "",
			expectedEnabled: true,
			expectedStyle:   StyleFull,
			expectedColors:  true,
		},
		{
			name:           "disabled banner",
			banner:         "false",
			style:          "",
			noColor:        "",
			expectedEnabled: false,
			expectedStyle:   StyleFull,
			expectedColors:  true,
		},
		{
			name:           "minimal style",
			banner:         "",
			style:          "minimal",
			noColor:        "",
			expectedEnabled: true,
			expectedStyle:   StyleMinimal,
			expectedColors:  true,
		},
		{
			name:           "no color",
			banner:         "",
			style:          "",
			noColor:        "1",
			expectedEnabled: true,
			expectedStyle:   StyleFull,
			expectedColors:  false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			_ = os.Setenv("GO_STARTER_BANNER", tt.banner)
			_ = os.Setenv("GO_STARTER_BANNER_STYLE", tt.style)
			_ = os.Setenv("NO_COLOR", tt.noColor)
			
			config := ConfigFromEnv()
			
			assert.Equal(t, tt.expectedEnabled, config.Enabled)
			assert.Equal(t, tt.expectedStyle, config.Style)
			assert.Equal(t, tt.expectedColors, config.Colors)
		})
	}
}

func TestGetBannerConfig(t *testing.T) {
	tests := []struct {
		name            string
		quiet           bool
		noBanner        bool
		bannerStyle     string
		expectedEnabled bool
		expectedQuiet   bool
		expectedStyle   BannerStyle
	}{
		{
			name:            "default config",
			quiet:           false,
			noBanner:        false,
			bannerStyle:     "",
			expectedEnabled: true,
			expectedQuiet:   false,
			expectedStyle:   StyleFull,
		},
		{
			name:            "quiet mode",
			quiet:           true,
			noBanner:        false,
			bannerStyle:     "",
			expectedEnabled: false,
			expectedQuiet:   true,
			expectedStyle:   StyleFull,
		},
		{
			name:            "no banner flag",
			quiet:           false,
			noBanner:        true,
			bannerStyle:     "",
			expectedEnabled: false,
			expectedQuiet:   false,
			expectedStyle:   StyleFull,
		},
		{
			name:            "minimal style",
			quiet:           false,
			noBanner:        false,
			bannerStyle:     "minimal",
			expectedEnabled: true,
			expectedQuiet:   false,
			expectedStyle:   StyleMinimal,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := GetBannerConfig(tt.quiet, tt.noBanner, tt.bannerStyle)
			
			assert.Equal(t, tt.expectedEnabled, config.Enabled)
			assert.Equal(t, tt.expectedQuiet, config.Quiet)
			assert.Equal(t, tt.expectedStyle, config.Style)
		})
	}
}

func TestRenderBanner(t *testing.T) {
	banner := RenderBanner("Test Banner")
	
	assert.NotEmpty(t, banner)
	assert.Contains(t, banner, "Test Banner")
}

func TestGopher(t *testing.T) {
	gopher := Gopher()
	
	assert.NotEmpty(t, gopher)
	// Gopher art contains characteristic ASCII art patterns
	assert.True(t, strings.Contains(gopher, "_") || strings.Contains(gopher, "*"))
}


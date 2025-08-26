// Package ascii provides professional ASCII art banners and logos for go-starter CLI.
// 
// This package offers multiple ASCII art styles optimized for terminal display:
//
// Main Banners:
//   - Banner(): Full professional banner with "GO-STARTER" text and enterprise messaging
//   - WelcomeBanner(): Welcome-focused banner with comprehensive feature messaging
//   - SlimBanner(): Single-line header suitable for documentation
//
// Logos and Compact Displays:
//   - Logo(): Medium-sized logo suitable for CLI startup
//   - CompactLogo(): Ultra-compact single-line logo for inline use
//   - MinimalBrand(): Clean boxed design for professional contexts
//
// All functions support:
//   - Color/monochrome modes via configuration
//   - Cross-platform terminal compatibility (Windows/macOS/Linux)
//   - Standard ASCII characters (no Unicode dependencies by default)
//   - Multiple display styles (full, minimal, none)
//   - Environment variable configuration (NO_COLOR, GO_STARTER_BANNER)
//
// Usage Examples:
//   fmt.Print(Banner())                    // Full professional banner
//   fmt.Print(CompactLogo())               // Inline: ⚡ go-starter
//   fmt.Print(MinimalBrand())              // Clean boxed design
//   fmt.Print(WelcomeBanner())             // Welcome screen with features
package ascii

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// BannerStyle represents different banner display styles
type BannerStyle int

const (
	StyleFull BannerStyle = iota
	StyleMinimal
	StyleNone
)

// BannerConfig controls banner display behavior
type BannerConfig struct {
	Enabled    bool
	Style      BannerStyle
	Colors     bool
	ShowOnHelp bool
	Quiet      bool
}

// DefaultConfig returns sensible default banner configuration
func DefaultConfig() *BannerConfig {
	return &BannerConfig{
		Enabled:    true,
		Style:      StyleFull,
		Colors:     true,
		ShowOnHelp: true,
		Quiet:      false,
	}
}

// ConfigFromEnv creates banner config from environment variables
func ConfigFromEnv() *BannerConfig {
	config := DefaultConfig()
	
	// Respect NO_COLOR standard
	if os.Getenv("NO_COLOR") != "" {
		config.Colors = false
	}
	
	// Check GO_STARTER_BANNER environment variable
	if banner := os.Getenv("GO_STARTER_BANNER"); banner != "" {
		config.Enabled = strings.ToLower(banner) != "false"
	}
	
	// Check GO_STARTER_BANNER_STYLE
	if style := os.Getenv("GO_STARTER_BANNER_STYLE"); style != "" {
		switch strings.ToLower(style) {
		case "minimal":
			config.Style = StyleMinimal
		case "none":
			config.Style = StyleNone
		case "full":
			config.Style = StyleFull
		}
	}
	
	return config
}

// Banner returns the full ASCII art banner for go-starter
func Banner() string {
	return BannerWithConfig(DefaultConfig())
}

// BannerWithConfig returns the ASCII art banner with custom configuration
func BannerWithConfig(config *BannerConfig) string {
	if !config.Enabled || config.Style == StyleNone {
		return ""
	}
	
	if config.Style == StyleMinimal {
		return LogoWithConfig(config)
	}
	
	// Define lipgloss styles
	var (
		cyanStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
		blueStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true) 
		greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
		yellowStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Bold(true)
	)
	
	// Disable colors if requested or not supported
	if !config.Colors || !termenv.HasDarkBackground() {
		cyanStyle = lipgloss.NewStyle()
		blueStyle = lipgloss.NewStyle()
		greenStyle = lipgloss.NewStyle()
		yellowStyle = lipgloss.NewStyle()
	}
	
	// Professional GO-STARTER banner with modern design
	banner := "\n" +
		cyanStyle.Render("   ██████╗  ██████╗       ███████╗████████╗ █████╗ ██████╗ ████████╗███████╗██████╗ ") + "\n" +
		cyanStyle.Render("  ██╔════╝ ██╔═══██╗      ██╔════╝╚══██╔══╝██╔══██╗██╔══██╗╚══██╔══╝██╔════╝██╔══██╗") + "\n" +
		blueStyle.Render("  ██║  ███╗██║   ██║█████╗███████╗   ██║   ███████║██████╔╝   ██║   █████╗  ██████╔╝") + "\n" +
		greenStyle.Render("  ██║   ██║██║   ██║╚════╝╚════██║   ██║   ██╔══██║██╔══██╗   ██║   ██╔══╝  ██╔══██╗") + "\n" +
		greenStyle.Render("  ╚██████╔╝╚██████╔╝      ███████║   ██║   ██║  ██║██║  ██║   ██║   ███████╗██║  ██║") + "\n" +
		greenStyle.Render("   ╚═════╝  ╚═════╝       ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝") + "\n" +
		"\n" +
		yellowStyle.Render("                        ⚡ Professional Go Project Generator ⚡") + "\n" +
		cyanStyle.Render("                            🏗️  12 Production-Ready Blueprints  🏗️") + "\n\n"
	
	return banner
}

// Logo returns a smaller ASCII logo for go-starter
func Logo() string {
	return LogoWithConfig(DefaultConfig())
}

// LogoWithConfig returns the ASCII logo with custom configuration
func LogoWithConfig(config *BannerConfig) string {
	if !config.Enabled || config.Style == StyleNone {
		return ""
	}
	
	var (
		cyanStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
		blueStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
		greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	)
	
	// Disable colors if requested
	if !config.Colors {
		cyanStyle = lipgloss.NewStyle().Bold(true)
		blueStyle = lipgloss.NewStyle().Bold(true)
		greenStyle = lipgloss.NewStyle().Bold(true)
	}
	
	// Compact professional logo
	return "\n" +
		cyanStyle.Render("   ██████╗  ██████╗      ███████╗") + "\n" +
		blueStyle.Render("  ██╔════╝ ██╔═══██╗     ██╔════╝") + "\n" +
		greenStyle.Render("  ██║  ███╗██║   ██║     ███████╗") + "\n" +
		greenStyle.Render("  ╚██████╔╝╚██████╔╝     ╚══════╝") + "\n" +
		"\n" +
		cyanStyle.Render("    🚀 go-starter - Enterprise Go Generator") + "\n\n"
}

// PrintLogo prints the logo to stdout with default configuration
func PrintLogo() {
	fmt.Print(Logo())
}

// PrintLogoWithConfig prints the logo with custom configuration
func PrintLogoWithConfig(config *BannerConfig) {
	fmt.Print(LogoWithConfig(config))
}

// PrintWelcome displays a welcome message with the banner
func PrintWelcome() {
	PrintWelcomeWithConfig(DefaultConfig())
}

// PrintWelcomeWithConfig displays a welcome message with custom configuration
func PrintWelcomeWithConfig(config *BannerConfig) {
	if config.Quiet || !config.Enabled {
		return
	}
	
	// Use the dedicated welcome banner for better messaging
	fmt.Print(WelcomeBannerWithConfig(config))
	
	// Additional context message
	contextStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true)
	
	if !config.Colors {
		contextStyle = lipgloss.NewStyle().Italic(true)
	}
	
	fmt.Println(contextStyle.Render("   Use 'go-starter new' to start your next Go project"))
	fmt.Println(contextStyle.Render("   For help: 'go-starter --help' or 'go-starter new --help'"))
	fmt.Println()
}

// RenderBanner renders custom text as a banner
func RenderBanner(text string) string {
	return RenderBannerWithConfig(text, DefaultConfig())
}

// RenderBannerWithConfig renders custom text as a banner with configuration
func RenderBannerWithConfig(text string, config *BannerConfig) string {
	if !config.Enabled || config.Style == StyleNone {
		return ""
	}
	
	// Simple banner rendering for custom text
	bannerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("8")).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)
	
	if !config.Colors {
		bannerStyle = bannerStyle.Foreground(lipgloss.NoColor{})
	}
	
	return bannerStyle.Render(text)
}

// Gopher returns a small ASCII gopher (kept for compatibility)
func Gopher() string {
	gopherStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	
	return gopherStyle.Render(`
        ,_---~~~~~----._         
 _,,_,*^____      _____*g*"*,   
/ __/ /'     ^.  /      \ ^@q   f 
[  @f | @))    |  | @))   l  0 _/  
 \ /   \~____ / __ \_____/    \   
  |           _l__l_           I   
  }          [______]           I  
  ]            | | |            |  
  ]             ~ ~             |  
  |                            |   
   |                           |
`)
}

// CompactLogo returns a very small logo for inline use
func CompactLogo() string {
	return CompactLogoWithConfig(DefaultConfig())
}

// CompactLogoWithConfig returns a compact logo with configuration
func CompactLogoWithConfig(config *BannerConfig) string {
	if !config.Enabled || config.Style == StyleNone {
		return ""
	}
	
	cyanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	
	if !config.Colors {
		cyanStyle = lipgloss.NewStyle().Bold(true)
	}
	
	// Ultra-compact single line logo
	return cyanStyle.Render("⚡ go-starter")
}

// SlimBanner returns a one-line banner for headers
func SlimBanner() string {
	return SlimBannerWithConfig(DefaultConfig())
}

// SlimBannerWithConfig returns a slim banner with configuration
func SlimBannerWithConfig(config *BannerConfig) string {
	if !config.Enabled || config.Style == StyleNone {
		return ""
	}
	
	var (
		cyanStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
		grayStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	)
	
	if !config.Colors {
		cyanStyle = lipgloss.NewStyle().Bold(true)
		grayStyle = lipgloss.NewStyle()
	}
	
	// Single line professional header
	return "\n" +
		cyanStyle.Render("███╗  ███ ██████╗     ████████╗████████╗ █████╗ ██████╗ ████████╗███████╗██████╗ ") + "\n" +
		grayStyle.Render("──────────────────────────────────────────────────────────────────────────────────────") + "\n" +
		cyanStyle.Render("                          ⚡ Professional Go Project Generator ⚡") + "\n\n"
}

// MinimalBrand returns a clean minimalist logo
func MinimalBrand() string {
	return MinimalBrandWithConfig(DefaultConfig())
}

// MinimalBrandWithConfig returns a minimalist brand with configuration
func MinimalBrandWithConfig(config *BannerConfig) string {
	if !config.Enabled || config.Style == StyleNone {
		return ""
	}
	
	var (
		cyanStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
		greenStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	)
	
	if !config.Colors {
		cyanStyle = lipgloss.NewStyle().Bold(true)
		greenStyle = lipgloss.NewStyle()
	}
	
	// Clean, modern minimalist design
	return "\n" +
		cyanStyle.Render("    ┌───────────────────────────────────────────────┐") + "\n" +
		cyanStyle.Render("    │                GO-STARTER                │") + "\n" +
		greenStyle.Render("    │          Enterprise Go Generator          │") + "\n" +
		cyanStyle.Render("    └───────────────────────────────────────────────┘") + "\n\n"
}

// WelcomeBanner returns a specialized welcome banner
func WelcomeBanner() string {
	return WelcomeBannerWithConfig(DefaultConfig())
}

// WelcomeBannerWithConfig returns a welcome banner with configuration
func WelcomeBannerWithConfig(config *BannerConfig) string {
	if !config.Enabled || config.Style == StyleNone {
		return ""
	}
	
	var (
		cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
		yellowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
		greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	)
	
	if !config.Colors {
		cyanStyle = lipgloss.NewStyle().Bold(true)
		yellowStyle = lipgloss.NewStyle()
		greenStyle = lipgloss.NewStyle()
	}
	
	// Welcome-focused banner with enterprise messaging
	return "\n" +
		cyanStyle.Render("██████████████████████████████████████████████████████████████████████████████████") + "\n" +
		cyanStyle.Render("██                        GO-STARTER                        ██") + "\n" +
		yellowStyle.Render("██              🏗️ Professional Go Project Generator 🏗️              ██") + "\n" +
		greenStyle.Render("██                  12 Production-Ready Blueprints                  ██") + "\n" +
		cyanStyle.Render("██                     Enterprise • Modern • Fast                     ██") + "\n" +
		cyanStyle.Render("██████████████████████████████████████████████████████████████████████████████████") + "\n\n"
}

// GetBannerConfig creates banner configuration from CLI flags
func GetBannerConfig(quiet bool, noBanner bool, bannerStyle string) *BannerConfig {
	config := ConfigFromEnv()
	
	// CLI flags override environment
	if quiet {
		config.Quiet = true
		config.Enabled = false
	}
	
	if noBanner {
		config.Enabled = false
	}
	
	if bannerStyle != "" {
		switch strings.ToLower(bannerStyle) {
		case "full":
			config.Style = StyleFull
		case "minimal":
			config.Style = StyleMinimal
		case "none":
			config.Style = StyleNone
		}
	}
	
	return config
}

// ShowAllVariants displays all available ASCII art variants (for testing/demo)
func ShowAllVariants() {
	config := DefaultConfig()
	
	fmt.Println("=== GO-STARTER ASCII ART VARIANTS ===")
	fmt.Println()
	
	fmt.Println("1. FULL BANNER (Default):")
	fmt.Print(BannerWithConfig(config))
	
	fmt.Println("2. STANDARD LOGO (Minimal Style):")
	fmt.Print(LogoWithConfig(config))
	
	fmt.Println("3. WELCOME BANNER:")
	fmt.Print(WelcomeBannerWithConfig(config))
	
	fmt.Println("4. SLIM BANNER (Header):")
	fmt.Print(SlimBannerWithConfig(config))
	
	fmt.Println("5. MINIMAL BRAND (Boxed):")
	fmt.Print(MinimalBrandWithConfig(config))
	
	fmt.Println("6. COMPACT LOGO (Inline):")
	fmt.Printf("   %s\n\n", CompactLogoWithConfig(config))
	
	fmt.Println("7. MONOCHROME VERSIONS (NO_COLOR=1):")
	monoConfig := &BannerConfig{
		Enabled: true,
		Style:   StyleFull,
		Colors:  false,
	}
	fmt.Print(BannerWithConfig(monoConfig))
	
	fmt.Println("=== END VARIANTS ===")
}
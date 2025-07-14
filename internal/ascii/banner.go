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
		cyanStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
		blueStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true) 
		greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Bold(true)
	)
	
	// Disable colors if requested or not supported
	if !config.Colors || !termenv.HasDarkBackground() {
		cyanStyle = lipgloss.NewStyle()
		blueStyle = lipgloss.NewStyle()
		greenStyle = lipgloss.NewStyle()
	}
	
	banner := "\n" +
		cyanStyle.Render("  ██████╗  ██████╗ ██╗      █████╗ ███╗   ██╗ ██████╗") + "\n" +
		cyanStyle.Render("  ██╔════╝ ██╔═══██╗██║     ██╔══██╗████╗  ██║██╔════╝") + "\n" +
		blueStyle.Render("  ██║  ███╗██║   ██║██║     ███████║██╔██╗ ██║██║  ███╗") + "\n" +
		blueStyle.Render("  ██║   ██║██║   ██║██║     ██╔══██║██║╚██╗██║██║   ██║") + "\n" +
		greenStyle.Render("  ╚██████╔╝╚██████╔╝███████╗██║  ██║██║ ╚████║╚██████╔╝") + "\n" +
		greenStyle.Render("   ╚═════╝  ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝") + "\n" +
		"\n" +
		cyanStyle.Render("  ███████╗████████╗ █████╗ ██████╗ ████████╗███████╗██████╗") + "\n" +
		blueStyle.Render("  ██╔════╝╚══██╔══╝██╔══██╗██╔══██╗╚══██╔══╝██╔════╝██╔══██╗") + "\n" +
		blueStyle.Render("  ███████╗   ██║   ███████║██████╔╝   ██║   █████╗  ██████╔╝") + "\n" +
		greenStyle.Render("  ╚════██║   ██║   ██╔══██║██╔══██╗   ██║   ██╔══╝  ██╔══██╗") + "\n" +
		greenStyle.Render("  ███████║   ██║   ██║  ██║██║  ██║   ██║   ███████╗██║  ██║") + "\n" +
		greenStyle.Render("  ╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚══════╝╚═╝  ╚═╝") + "\n\n"
	
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
	
	cyanStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Bold(true)
	
	// Disable colors if requested
	if !config.Colors {
		cyanStyle = lipgloss.NewStyle().Bold(true)
	}
	
	return cyanStyle.Render(`
 _____ _____ 
|   __|     |
|  |  |  |  |
|_____|_____|`) + " " + cyanStyle.Render("STARTER") + "\n"
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
	
	fmt.Print(BannerWithConfig(config))
	
	// Welcome message styling
	welcomeStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Bold(true).
		MarginTop(1)
	
	if !config.Colors {
		welcomeStyle = lipgloss.NewStyle().Bold(true)
	}
	
	fmt.Println(welcomeStyle.Render("🚀 Welcome to Go-Starter!"))
	fmt.Println(welcomeStyle.Render("   Generate Go projects with modern best practices"))
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